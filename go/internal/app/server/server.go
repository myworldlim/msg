// internal/app/server/server.go
package server

import (
	"chitchat/internal/middleware/cors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(db *pgxpool.Pool) *mux.Router {
	r := mux.NewRouter()

	// Глобальные middleware - ПОРЯДОК ВАЖЕН!
	r.Use(loggingMiddleware)  // 1. Логирование запросов
	r.Use(recoveryMiddleware) // 2. Восстановление после паники
	// Rate limiting removed — middleware is a no-op or omitted
	r.Use(cors.CORSMiddleware()) // 4. CORS

	// Настройка маршрутов
	SetupRoutes(r, db)

	return r
}

// loggingMiddleware логирует запросы
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// recoveryMiddleware восстанавливает после паники
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
