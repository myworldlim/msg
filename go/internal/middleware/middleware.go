// internal/middleware/middleware.go
package middleware

import (
	"chitchat/internal/middleware/cors"
	"net/http"
)

// CORSMiddleware прокси-функция для CORS
func CORSMiddleware() func(http.Handler) http.Handler {
	return cors.CORSMiddleware()
}