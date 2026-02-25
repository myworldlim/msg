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

// LoginPasswordRequest DTO для входа по паролю
type LoginPasswordRequest struct {
	UserUid  string `json:"userUid"`
	Password string `json:"password"`
}

// LoginPasswordResponse DTO для ответа при входе по паролю
type LoginPasswordResponse struct {
	Success    bool `json:"success"`
	Message    string `json:"message"`
	Protection bool `json:"protection"`
}

// LoginPasswordHandler обрабатывает вход по паролю
// POST /auth/password/login
func LoginPasswordHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}

		// Валидация
		if req.UserUid == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "userUid and password are required"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Проверяем блокировку перед попыткой входа
		errorActive, _, lockedUntil, err := repo.CheckPasswordError(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (check password error): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Если аккаунт заблокирован
		if errorActive && lockedUntil != nil && time.Now().Before(*lockedUntil) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Account temporarily locked",
				"locked_until": lockedUntil.Unix(),
			})
			return
		}

		// Получаем хеш пароля из БД
		storedHash, found, err := repo.GetPasswordHashByUserUID(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (get password hash): %v", err)
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
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid stored password format"})
			return
		}

		// Первые 16 байт - соль, остальное - хеш
		salt := storedBytes[:16]
		storedHashBytes := storedBytes[16:]

		// Хешируем введенный пароль с той же солью
		inputHash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)

		// Сравниваем хеши
		if !equalBytes(inputHash, storedHashBytes) {
			// Неправильный пароль - увеличиваем счетчик ошибок
			if err := repo.IncrementPasswordError(ctx, db, req.UserUid); err != nil {
				log.Printf("DB error (increment password error): %v", err)
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}

		// Пароль правильный - сбрасываем счетчик ошибок
		if err := repo.ResetPasswordError(ctx, db, req.UserUid); err != nil {
			log.Printf("DB error (reset password error): %v", err)
			// Не прерываем выполнение, просто логируем
		}

		// Получаем статус protection
		protectionStatus, err := repo.GetProtectionStatusByUserUID(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (get protection status): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Если дополнительная защита отключена - создаем сессию
		if !protectionStatus {
			// Находим guid_id для создания сессии
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
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginPasswordResponse{
			Success:    true,
			Message:    "Login successful",
			Protection: protectionStatus,
		})
	}
}

// equalBytes безопасно сравнивает два байтовых массива
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}