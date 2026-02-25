package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	repo "chitchat/internal/repository/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrorPasswordRequest DTO для запроса
type ErrorPasswordRequest struct {
	UserUid string `json:"userUid"`
}

// ErrorPasswordResponse DTO для ответа
type ErrorPasswordResponse struct {
	Success        bool   `json:"success"`
	ErrorActive    bool   `json:"error_active"`
	FailedAttempts int    `json:"failed_attempts"`
	LockedUntil    *time.Time `json:"locked_until,omitempty"`
	TimeRemaining  int    `json:"time_remaining"`
	Error          string `json:"error,omitempty"`
}

// CheckErrorPasswordHandler проверяет статус блокировки пароля
// POST /auth/password/error
// Body: { "userUid": "..." }
func CheckErrorPasswordHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// Парсим JSON request
		var req ErrorPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorPasswordResponse{
				Success: false,
				Error:   "Invalid JSON",
			})
			return
		}

		// Валидируем userUid
		if req.UserUid == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorPasswordResponse{
				Success: false,
				Error:   "userUid is required",
			})
			return
		}

		// Проверяем статус блокировки в БД
		errorActive, failedAttempts, lockedUntil, err := repo.CheckPasswordError(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (check password error): %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorPasswordResponse{
				Success: false,
				Error:   "Internal server error",
			})
			return
		}

		// Вычисляем оставшееся время блокировки
		var timeRemaining int
		if lockedUntil != nil && time.Now().Before(*lockedUntil) {
			timeRemaining = int(lockedUntil.Sub(time.Now()).Seconds())
		}

		// Возвращаем результат
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		response := ErrorPasswordResponse{
			Success:        true,
			ErrorActive:    errorActive,
			FailedAttempts: failedAttempts,
			TimeRemaining:  timeRemaining,
		}
		
		if lockedUntil != nil {
			response.LockedUntil = lockedUntil
		}
		
		json.NewEncoder(w).Encode(response)
	}
}