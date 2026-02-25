package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateSession inserts a new session record
func CreateSession(ctx context.Context, db *pgxpool.Pool, userUid string, guidID int64, sessionToken string, sessionRefresh string, userAgent string, ip string, sessionExpires time.Time, refreshExpires time.Time) error {
	_, err := db.Exec(ctx,
		`INSERT INTO sessions (user_uid, guid_id, session_token, session_refresh, session_user_agent, ip_address, created_at, session_expires_at, session_refresh_expires_at)
         VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, $7, $8)`,
		userUid, guidID, sessionToken, sessionRefresh, userAgent, ip, sessionExpires, refreshExpires)
	return err
}

// GetSessionByToken returns session info by session_token
func GetSessionByToken(ctx context.Context, db *pgxpool.Pool, token string) (map[string]interface{}, bool, error) {
	row := db.QueryRow(ctx, `SELECT session_id, user_uid, guid_id, session_token, session_refresh, session_user_agent, ip_address, created_at, session_expires_at, session_refresh_expires_at FROM sessions WHERE session_token = $1 LIMIT 1`, token)
	var sessionID int64
	var userUid string
	var guidID *int64
	var sessionToken string
	var sessionRefresh string
	var userAgent, ip string
	var createdAt, sessExp, refreshExp *time.Time

	err := row.Scan(&sessionID, &userUid, &guidID, &sessionToken, &sessionRefresh, &userAgent, &ip, &createdAt, &sessExp, &refreshExp)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	res := map[string]interface{}{
		"session_id":                 sessionID,
		"user_uid":                   userUid,
		"guid_id":                    guidID,
		"session_token":              sessionToken,
		"session_refresh":            sessionRefresh,
		"session_user_agent":         userAgent,
		"ip_address":                 ip,
		"created_at":                 createdAt,
		"session_expires_at":         sessExp,
		"session_refresh_expires_at": refreshExp,
	}
	return res, true, nil
}

// GetSessionByRefresh finds session by refresh token
func GetSessionByRefresh(ctx context.Context, db *pgxpool.Pool, refresh string) (map[string]interface{}, bool, error) {
	row := db.QueryRow(ctx, `SELECT session_id, user_uid, guid_id, session_token, session_refresh, session_user_agent, ip_address, created_at, session_expires_at, session_refresh_expires_at FROM sessions WHERE session_refresh = $1 LIMIT 1`, refresh)
	var sessionID int64
	var userUid string
	var guidID *int64
	var sessionToken string
	var sessionRefresh string
	var userAgent, ip string
	var createdAt, sessExp, refreshExp *time.Time

	err := row.Scan(&sessionID, &userUid, &guidID, &sessionToken, &sessionRefresh, &userAgent, &ip, &createdAt, &sessExp, &refreshExp)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	res := map[string]interface{}{
		"session_id":                 sessionID,
		"user_uid":                   userUid,
		"guid_id":                    guidID,
		"session_token":              sessionToken,
		"session_refresh":            sessionRefresh,
		"session_user_agent":         userAgent,
		"ip_address":                 ip,
		"created_at":                 createdAt,
		"session_expires_at":         sessExp,
		"session_refresh_expires_at": refreshExp,
	}
	return res, true, nil
}

// UpdateSessionTokenByRefresh rotates the session_token when refresh is valid
func UpdateSessionTokenByRefresh(ctx context.Context, db *pgxpool.Pool, refresh string, newToken string, newExp time.Time) error {
	_, err := db.Exec(ctx, `UPDATE sessions SET session_token = $1, session_expires_at = $2 WHERE session_refresh = $3 AND session_refresh_expires_at > CURRENT_TIMESTAMP`, newToken, newExp, refresh)
	return err
}

// DeleteSessionByToken удаляет сессию по токену
func DeleteSessionByToken(ctx context.Context, db *pgxpool.Pool, token string) error {
	_, err := db.Exec(ctx, `DELETE FROM sessions WHERE session_token = $1`, token)
	return err
}

// DeleteSessionByRefresh удаляет сессию по refresh токену
func DeleteSessionByRefresh(ctx context.Context, db *pgxpool.Pool, refresh string) error {
	_, err := db.Exec(ctx, `DELETE FROM sessions WHERE session_refresh = $1`, refresh)
	return err
}

// CheckUserBlocked проверяет заблокирован ли пользователь по user_uid
// Возвращает (blocked bool, reason string, error)
func CheckUserBlocked(ctx context.Context, db *pgxpool.Pool, userUid string) (bool, string, error) {
	var guidID int64

	// 1. Найти guid_id по user_uid
	err := db.QueryRow(ctx, `SELECT guid_id FROM guid WHERE user_uid = $1`, userUid).Scan(&guidID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Пользователь не найден = не заблокирован
			return false, "", nil
		}
		return false, "", err
	}

	// 2. Проверить в таблице blocked
	var blockedStatus bool
	var blockedType string

	err = db.QueryRow(ctx,
		`SELECT blocked_status, blocked_type FROM blocked WHERE guid_id = $1 AND blocked_status = true LIMIT 1`,
		guidID,
	).Scan(&blockedStatus, &blockedType)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Запись о блокировке не найдена = не заблокирован
			return false, "", nil
		}
		return false, "", err
	}

	// Если blocked_status = true, вернуть true и reason (из blocked_type)
	if blockedStatus {
		return true, blockedType, nil
	}

	return false, "", nil
}
