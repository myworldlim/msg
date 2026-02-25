package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	repo "chitchat/internal/repository/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"
)

// CreateSecretRequest DTO для создания секретного слова
type CreateSecretRequest struct {
	UserUid    string `json:"userUid"`
	SecretWord string `json:"secretWord"`
}

// CreateSecretResponse DTO для ответа при создании секретного слова
type CreateSecretResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateSecretHandler обрабатывает создание секретного слова
// POST /auth/secret/create
func CreateSecretHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSecretRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}

		// Валидация
		if req.UserUid == "" || req.SecretWord == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "userUid and secretWord are required"})
			return
		}

		if len(req.SecretWord) < 3 || len(req.SecretWord) > 100 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Secret word must be between 3 and 100 characters"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Находим guid_id по user_uid
		guidID, found, err := repo.FindGuidIDByUserUID(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (find guid): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}

		// Хешируем секретное слово с рандомной солью
		salt := make([]byte, 16)
		_, _ = rand.Read(salt)
		hash := argon2.IDKey([]byte(req.SecretWord), salt, 1, 64*1024, 4, 32)
		// Сохраняем соль + хеш вместе
		hashWithSalt := append(salt, hash...)
		hashHex := hex.EncodeToString(hashWithSalt)

		// Сохраняем в БД
		err = repo.CreateSecret(ctx, db, guidID, hashHex)
		if err != nil {
			log.Printf("DB error (create secret): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create secret"})
			return
		}

		// Создаем сессию и cookies (пользователь полностью авторизован)
		tokenBytes := make([]byte, 32)
		_, _ = rand.Read(tokenBytes)
		sessionToken := hex.EncodeToString(tokenBytes)

		refreshBytes := make([]byte, 48)
		_, _ = rand.Read(refreshBytes)
		sessionRefresh := hex.EncodeToString(refreshBytes)

		sessionExp := time.Now().Add(24 * time.Hour)
		refreshExp := time.Now().Add(30 * 24 * time.Hour)

		userAgent := r.Header.Get("User-Agent")
		ip := r.RemoteAddr

		// Сохраняем сессию в БД
		if err := repo.CreateSession(ctx, db, req.UserUid, guidID, sessionToken, sessionRefresh, userAgent, ip, sessionExp, refreshExp); err != nil {
			log.Printf("DB error (create session): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create session"})
			return
		}

		// Устанавливаем cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  sessionExp,
			SameSite: http.SameSiteLaxMode,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "session_refresh",
			Value:    sessionRefresh,
			Path:     "/",
			HttpOnly: true,
			Expires:  refreshExp,
			SameSite: http.SameSiteLaxMode,
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateSecretResponse{
			Success: true,
			Message: "Secret word created successfully",
		})
	}
}