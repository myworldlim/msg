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

// BlockedStatusRequest DTO для запроса
type BlockedStatusRequest struct {
	UserUid string `json:"userUid"`
}

// BlockedStatusResponse DTO для ответа
type BlockedStatusResponse struct {
	Success bool   `json:"success"`
	Blocked bool   `json:"blocked"`
	UserUid string `json:"userUid,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}

// CheckBlockedHandler проверяет заблокирован ли пользователь
// POST /auth/blocked
// Body: { "userUid": "..." }
func CheckBlockedHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// Парсим JSON request
		var req BlockedStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(BlockedStatusResponse{
				Success: false,
				Blocked: false,
				Error:   "Invalid JSON",
			})
			return
		}

		// Валидируем userUid
		if req.UserUid == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(BlockedStatusResponse{
				Success: false,
				Blocked: false,
				Error:   "userUid is required",
			})
			return
		}

		// Проверяем блокировку в БД
		blocked, reason, err := repo.CheckUserBlocked(ctx, db, req.UserUid)
		if err != nil {
			log.Printf("DB error (check blocked): %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(BlockedStatusResponse{
				Success: false,
				Blocked: false,
				UserUid: req.UserUid,
				Error:   "Internal server error",
			})
			return
		}

		// Возвращаем результат
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BlockedStatusResponse{
			Success: true,
			Blocked: blocked,
			UserUid: req.UserUid,
			Reason:  reason,
		})
	}
}
