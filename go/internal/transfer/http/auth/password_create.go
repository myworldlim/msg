package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	repo "chitchat/internal/repository/auth"
	"crypto/rand"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"
)

// RegisterPasswordRequest DTO для регистрации пароля
type RegisterPasswordRequest struct {
	UserUid    string `json:"userUid"`
	Password   string `json:"password"`
	Protection bool   `json:"protection"`
}

// RegisterPasswordResponse DTO для ответа при регистрации пароля
type RegisterPasswordResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Protection bool   `json:"protection"`
}

// RegisterPasswordHandler обрабатывает регистрацию пароля пользователя
// POST /auth/password/register
func RegisterPasswordHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}

		// Валидация входных данных
		if req.UserUid == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "userUid and password are required"})
			return
		}

		if len(req.Password) < 8 || len(req.Password) > 128 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Password must be between 8 and 128 characters"})
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

		// ✅ Проверка: если пароль уже активен, не разрешаем повторную регистрацию
		isPasswordActive, err := repo.CheckPasswordActiveByGuid(ctx, db, guidID)
		if err != nil {
			log.Printf("DB error (check password active): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		if isPasswordActive {
			w.WriteHeader(http.StatusConflict) // 409
			json.NewEncoder(w).Encode(map[string]string{"error": "Password already registered for this user"})
			return
		}

		// Хешируем пароль через argon2id с рандомной солью
		salt := make([]byte, 16)
		_, _ = rand.Read(salt)
		hash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)
		// Сохраняем соль + хеш вместе
		hashWithSalt := append(salt, hash...)
		hashHex := hex.EncodeToString(hashWithSalt)

		// Сохраняем пароль в БД (обновляем запись password с password_active = true)
		err = repo.CreatePassword(ctx, db, guidID, hashHex)
		if err != nil {
			log.Printf("DB error (create password): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register password"})
			return
		}

		// Сохраняем protection статус
		err = repo.CreateProtection(ctx, db, guidID, req.Protection)
		if err != nil {
			log.Printf("DB error (create protection): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save protection status"})
			return
		}

		// Создаем запись в secrets с secret_active = false (для всех пользователей)
		err = repo.CreateEmptySecret(ctx, db, guidID)
		if err != nil {
			log.Printf("DB error (create empty secret): %v", err)
			// Не прерываем выполнение, это не критично
		}

		// Если дополнительная защита не включена — создаём сессию и ставим HttpOnly cookies
		if !req.Protection {
			// generate tokens
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

			// сохраняем в БД
			if err := repo.CreateSession(ctx, db, req.UserUid, guidID, sessionToken, sessionRefresh, userAgent, ip, sessionExp, refreshExp); err != nil {
				log.Printf("DB error (create session): %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create session"})
				return
			}

			// set cookies
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
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RegisterPasswordResponse{
			Success:    true,
			Message:    "Password registered successfully",
			Protection: req.Protection,
		})
	}
}