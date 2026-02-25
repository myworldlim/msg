// internal/app/server/routes.go
package server

import (
	"encoding/json"
	"net/http"

	authhttp "chitchat/internal/transfer/http/auth"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *mux.Router, db *pgxpool.Pool) {
	// Глобальный обработчик OPTIONS для preflight запросов (CORS)
	r.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
	}).Methods("OPTIONS")

	// auth routes group
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/open", authhttp.OpenHandler(db)).Methods("POST")
	authRouter.HandleFunc("/password/register", authhttp.RegisterPasswordHandler(db)).Methods("POST")
	authRouter.HandleFunc("/session/check", authhttp.SessionCheckHandler(db)).Methods("GET")
	authRouter.HandleFunc("/session/refresh", authhttp.SessionRefreshHandler(db)).Methods("POST")
	authRouter.HandleFunc("/logout", authhttp.LogoutHandler(db)).Methods("POST")
	authRouter.HandleFunc("/secret/create", authhttp.CreateSecretHandler(db)).Methods("POST")
	authRouter.HandleFunc("/password/login", authhttp.LoginPasswordHandler(db)).Methods("POST")
	authRouter.HandleFunc("/password/error", authhttp.CheckErrorPasswordHandler(db)).Methods("POST")
	authRouter.HandleFunc("/password/recover", authhttp.CheckRecoverPasswordHandler(db)).Methods("POST")
	authRouter.HandleFunc("/secret/login", authhttp.LoginSecretHandler(db)).Methods("POST")
	authRouter.HandleFunc("/blocked", authhttp.CheckBlockedHandler(db)).Methods("POST")

	// root — простой статусный endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие cookies: "session_token" и "session_refresh"
		_, errToken := r.Cookie("session_token")
		_, errRefresh := r.Cookie("session_refresh")

		tokenExists := errToken == nil
		refreshExists := errRefresh == nil

		resp := map[string]interface{}{
			"session_token":   tokenExists,
			"session_refresh": refreshExists,
			"has_session":     tokenExists && refreshExists,
		}

		w.Header().Set("Content-Type", "application/json")
		if tokenExists && refreshExists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")
}
