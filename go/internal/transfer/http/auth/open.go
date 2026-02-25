package auth

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	repo "chitchat/internal/repository/auth"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Request and response DTOs
type OpenRequest struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

type OpenResponse struct {
	IsUserUid    string `json:"isUserUid"`
	IsExists     bool   `json:"isExists"`
	IsPassword   bool   `json:"isPassword"`
	IsBlocked    bool   `json:"isBlocked"`
	IsIdentifier string `json:"isIdentifier"`
	IsUserAgent  string `json:"isUserAgent"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Handler for POST /auth/open
func OpenHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OpenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		req.Identifier = strings.TrimSpace(req.Identifier)
		if req.Identifier == "" || (req.Type != "email" && req.Type != "number") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "identifier and valid type are required"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// validate input
		if req.Type == "email" {
			if !isValidEmail(req.Identifier) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email format"})
				return
			}
			req.Identifier = strings.ToLower(req.Identifier)
		} else {
			// normalize phone
			n := normalizePhone(req.Identifier)
			if len(n) < 8 || len(n) > 15 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid phone number"})
				return
			}
			req.Identifier = n
		}

		// Extract user agent from request
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			userAgent = "Unknown"
		}

		// repository lookups
		var uid string
		var exists bool
		var isBlocked bool
		var err error

		if req.Type == "email" {
			uid, exists, err = repo.FindUserUIDByEmail(ctx, db, req.Identifier)
		} else {
			uid, exists, err = repo.FindUserUIDByNumber(ctx, db, req.Identifier)
		}
		if err != nil {
			log.Printf("DB error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}

		if !exists {
			// user not found — create new user in database
			// Generate a random user_uid (20 chars alphanumeric)
			newUID := generateUserUID()

			var email, number string
			if req.Type == "email" {
				email = req.Identifier
			} else {
				number = req.Identifier
			}

			_, err := repo.CreateOrGetUser(ctx, db, newUID, email, number)
			if err != nil {
				log.Printf("DB error (create user): %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
				return
			}
			uid = newUID
			isBlocked = false

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(OpenResponse{
				IsUserUid:    uid,
				IsExists:     false,
				IsPassword:   false, // новый пользователь — пароль еще не установлен
				IsBlocked:    false,
				IsIdentifier: req.Identifier,
				IsUserAgent:  userAgent,
			})
			return
		}

		// found user — check guid
		guidID, foundGuid, err := repo.FindGuidIDByUserUID(ctx, db, uid)
		if err != nil {
			log.Printf("DB error (guid): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}
		if !foundGuid {
			// no guid record — consider not blocked and no password
			isBlocked = false
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(OpenResponse{
				IsUserUid:    uid,
				IsExists:     true,
				IsPassword:   false,
				IsBlocked:    false,
				IsIdentifier: req.Identifier,
				IsUserAgent:  userAgent,
			})
			return
		}

		// Check password status
		hasPassword, err := repo.CheckPasswordActiveByGuid(ctx, db, guidID)
		if err != nil {
			log.Printf("DB error (password check): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}

		isBlocked, err = repo.GetBlockedStatusByGuidID(ctx, db, guidID)
		if err != nil {
			log.Printf("DB error (blocked): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OpenResponse{
			IsUserUid:    uid,
			IsExists:     true,
			IsPassword:   hasPassword, // true если пароль активирован
			IsBlocked:    isBlocked,
			IsIdentifier: req.Identifier,
			IsUserAgent:  userAgent,
		})
	}
}

// basic email regex
var emailRegexp = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func isValidEmail(email string) bool {
	return emailRegexp.MatchString(email)
}

// normalizePhone removes non-digit characters except leading +
func normalizePhone(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	// keep leading + if present
	keepPlus := false
	if strings.HasPrefix(s, "+") {
		keepPlus = true
	}
	var b strings.Builder
	if keepPlus {
		b.WriteByte('+')
	}
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// generateUserUID generates a random 20-character alphanumeric user UID
func generateUserUID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	uid := make([]byte, 20)
	for i := range uid {
		uid[i] = charset[rand.Intn(len(charset))]
	}
	return string(uid)
}
