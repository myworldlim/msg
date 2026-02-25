package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateSecret создает секретное слово для пользователя
func CreateSecret(ctx context.Context, db *pgxpool.Pool, guidID int64, secretHash string) error {
	_, err := db.Exec(ctx,
		`INSERT INTO secrets (guid_id, secret_hash, secret_active, created_at) 
		 VALUES ($1, $2, true, $3)
		 ON CONFLICT (guid_id) DO UPDATE SET 
		 secret_hash = EXCLUDED.secret_hash, 
		 secret_active = true, 
		 created_at = EXCLUDED.created_at`,
		guidID, secretHash, time.Now())
	return err
}

// CheckSecretActiveByGuid проверяет, есть ли активное секретное слово
func CheckSecretActiveByGuid(ctx context.Context, db *pgxpool.Pool, guidID int64) (bool, error) {
	var exists bool
	row := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM secrets WHERE guid_id = $1 AND secret_active = true)`, guidID)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetSecretHashByGuid получает хеш секретного слова
func GetSecretHashByGuid(ctx context.Context, db *pgxpool.Pool, guidID int64) (string, bool, error) {
	var hash string
	row := db.QueryRow(ctx, `SELECT secret_hash FROM secrets WHERE guid_id = $1 AND secret_active = true LIMIT 1`, guidID)
	err := row.Scan(&hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return hash, true, nil
}

// CreateEmptySecret создает пустую запись в secrets с secret_active = false
func CreateEmptySecret(ctx context.Context, db *pgxpool.Pool, guidID int64) error {
	_, err := db.Exec(ctx,
		`INSERT INTO secrets (guid_id, secret_hash, secret_active, created_at) 
		 VALUES ($1, '', false, $2)
		 ON CONFLICT (guid_id) DO NOTHING`, // не обновляем если уже есть
		guidID, time.Now())
	return err
}