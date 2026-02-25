package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetPasswordHashByUserUID получает хеш пароля по user_uid
func GetPasswordHashByUserUID(ctx context.Context, db *pgxpool.Pool, userUID string) (string, bool, error) {
	var hash string
	row := db.QueryRow(ctx, 
		`SELECT p.password_hash 
		 FROM passwords p 
		 JOIN guid g ON p.guid_id = g.guid_id 
		 WHERE g.user_uid = $1 AND p.password_active = true 
		 LIMIT 1`, userUID)
	
	err := row.Scan(&hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return hash, true, nil
}

// GetProtectionStatusByUserUID получает статус protection по user_uid
func GetProtectionStatusByUserUID(ctx context.Context, db *pgxpool.Pool, userUID string) (bool, error) {
	var status bool
	row := db.QueryRow(ctx,
		`SELECT COALESCE(pr.protection_status, false)
		 FROM protection pr
		 JOIN guid g ON pr.guid_id = g.guid_id
		 WHERE g.user_uid = $1
		 LIMIT 1`, userUID)
	
	err := row.Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return status, nil
}