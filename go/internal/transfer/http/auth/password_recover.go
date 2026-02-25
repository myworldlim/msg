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

// RecoverPasswordRequest DTO для запроса
type RecoverPasswordRequest struct {
	UserUid string `json:"userUid"`
}

// RecoverPasswordResponse DTO для ответа
type RecoverPasswordResponse struct {
	Success           bool   `json:"success"`
	RecoveryAvailable bool   `json:"recovery_available"`
	RecoveryMethod    string `json:"recovery_method,omitempty"`
	RecoveryContact   string `json:"recovery_contact,omitempty"`
	Error             string `json:"error,omitempty"`
}

// CheckRecoverPasswordHandler проверяет доступность восстановления пароля
// POST /auth/password/recover
// Body: { "userUid": "KtabQdxlMHUIcq4VU8d9" }
func CheckRecoverPasswordHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// Парсим JSON request
		var req RecoverPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RecoverPasswordResponse{
				Success: false,
				Error:   "Invalid JSON",
			})
			return
		}

		// Валидируем userUid
		if req.UserUid == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RecoverPasswordResponse{
				Success: false,
				Error:   "userUid is required",
			})
			return
		}

		// Проверяем доступность восстановления в БД
		available, method, contact, err := repo.CheckPasswordRecovery(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (check password recovery): %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(RecoverPasswordResponse{
				Success: false,
				Error:   "Internal server error",
			})
			return
		}

		// Возвращаем результат
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RecoverPasswordResponse{
			Success:           true,
			RecoveryAvailable: available,
			RecoveryMethod:    method,
			RecoveryContact:   contact,
		})
	}
}