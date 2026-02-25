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
)

// SessionCheckResponse DTO
type SessionCheckResponse struct {
	HasSession bool   `json:"hasSession"`
	UserUid    string `json:"userUid,omitempty"`
}

// SessionCheckHandler checks session_token cookie and refreshes if needed
func SessionCheckHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// try session_token
		tokenCookie, err := r.Cookie("session_token")
		if err == nil {
			session, found, err := repo.GetSessionByToken(ctx, db, tokenCookie.Value)
			if err != nil {
				log.Printf("DB error (get session by token): %v", err)
			} else if found {
				// check expiry
				if sessExp, ok := session["session_expires_at"].(*time.Time); ok && sessExp != nil && sessExp.After(time.Now()) {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(SessionCheckResponse{HasSession: true, UserUid: session["user_uid"].(string)})
					return
				}
			}
		}

		// try refresh token
		refreshCookie, err := r.Cookie("session_refresh")
		if err == nil {
			session, found, err := repo.GetSessionByRefresh(ctx, db, refreshCookie.Value)
			if err != nil {
				log.Printf("DB error (get session by refresh): %v", err)
			} else if found {
				// check refresh expiry
				if refreshExp, ok := session["session_refresh_expires_at"].(*time.Time); ok && refreshExp != nil && refreshExp.After(time.Now()) {
					// rotate token
					newTokenBytes := make([]byte, 32)
					_, _ = rand.Read(newTokenBytes)
					newToken := hex.EncodeToString(newTokenBytes)
					newExp := time.Now().Add(24 * time.Hour)
					if err := repo.UpdateSessionTokenByRefresh(ctx, db, refreshCookie.Value, newToken, newExp); err != nil {
						log.Printf("DB error (update session token): %v", err)
					} else {
						// set cookie
						http.SetCookie(w, &http.Cookie{
							Name:     "session_token",
							Value:    newToken,
							Path:     "/",
							HttpOnly: true,
							Expires:  newExp,
							SameSite: http.SameSiteLaxMode,
						})
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(SessionCheckResponse{HasSession: true, UserUid: session["user_uid"].(string)})
						return
					}
				}
			}
		}

		// no valid session
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(SessionCheckResponse{HasSession: false})
	}
}

// SessionRefreshHandler — rotate refresh (optional endpoint)
func SessionRefreshHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		refreshCookie, err := r.Cookie("session_refresh")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		session, found, err := repo.GetSessionByRefresh(ctx, db, refreshCookie.Value)
		if err != nil || !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// check expiry
		if refreshExp, ok := session["session_refresh_expires_at"].(*time.Time); !ok || refreshExp == nil || refreshExp.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// generate new tokens
		newTokenBytes := make([]byte, 32)
		_, _ = rand.Read(newTokenBytes)
		newToken := hex.EncodeToString(newTokenBytes)
		newExp := time.Now().Add(24 * time.Hour)
		if err := repo.UpdateSessionTokenByRefresh(ctx, db, refreshCookie.Value, newToken, newExp); err != nil {
			log.Printf("DB error (update session token): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    newToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  newExp,
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusOK)
	}
}

// LogoutHandler удаляет сессию и очищает cookies
func LogoutHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// Получаем токены из cookies
		tokenCookie, _ := r.Cookie("session_token")
		refreshCookie, _ := r.Cookie("session_refresh")

		// Удаляем сессию из БД
		if tokenCookie != nil {
			repo.DeleteSessionByToken(ctx, db, tokenCookie.Value)
		}
		if refreshCookie != nil {
			repo.DeleteSessionByRefresh(ctx, db, refreshCookie.Value)
		}

		// Очищаем cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "session_refresh",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
