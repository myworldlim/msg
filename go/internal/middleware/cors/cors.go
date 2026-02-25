// internal/middleware/cors/cors.go
package cors

import (
	"chitchat/config"
	"net/http"
	"strings"
)

// CORSMiddleware создает middleware для CORS
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Проверяем, разрешен ли источник
			allowed := isOriginAllowed(origin)

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 часа
			}

			// Обработка preflight запросов
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string) bool {
	// В dev-окружении разрешаем все origins
	if config.AppConfig.AppEnv == "development" || config.AppConfig.AppEnv == "local" {
		return true
	}

	// Если нет настроек CORS, разрешаем все (для разработки)
	if len(config.AppConfig.CORSAccepted) == 0 {
		return true
	}

	// Проверяем каждый разрешенный источник
	for _, allowedOrigin := range config.AppConfig.CORSAccepted {
		allowedOrigin = strings.TrimSpace(allowedOrigin)
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}

	return false
}
