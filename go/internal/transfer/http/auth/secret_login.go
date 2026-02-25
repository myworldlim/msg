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

// LoginSecretRequest DTO для входа по секретному слову
type LoginSecretRequest struct {
	UserUid    string `json:"userUid"`
	SecretWord string `json:"secretWord"`
}

// LoginSecretResponse DTO для ответа при входе по секретному слову
type LoginSecretResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LoginSecretHandler обрабатывает вход по секретному слову
// POST /auth/secret/login
func LoginSecretHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginSecretRequest
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

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Получаем хеш секретного слова из БД
		storedHash, found, err := repo.GetSecretHashByUserUID(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (get secret hash): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}

		// Распаковываем соль и хеш из сохраненного значения
		storedBytes, err := hex.DecodeString(storedHash)
		if err != nil || len(storedBytes) < 16 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid stored secret format"})
			return
		}

		// Первые 16 байт - соль, остальное - хеш
		salt := storedBytes[:16]
		storedHashBytes := storedBytes[16:]

		// Хешируем введенное секретное слово с той же солью
		inputHash := argon2.IDKey([]byte(req.SecretWord), salt, 1, 64*1024, 4, 32)

		// Сравниваем хеши
		if !equalBytes(inputHash, storedHashBytes) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}

		// Секретное слово верное - создаем сессию
		guidID, found, err := repo.FindGuidIDByUserUID(ctx, db, req.UserUid)
		if err != nil || !found {
			log.Printf("DB error (find guid): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Создаем токены сессии
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
		json.NewEncoder(w).Encode(LoginSecretResponse{
			Success: true,
			Message: "Login successful",
		})
	}
}