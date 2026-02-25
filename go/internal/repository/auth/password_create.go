package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CheckPasswordActiveByGuid проверяет, активен ли пароль (password_active) для guid_id
// Возвращает true если пароль активирован, false если нет или если записи нет
func CheckPasswordActiveByGuid(ctx context.Context, db *pgxpool.Pool, guidID int64) (bool, error) {
	var active bool
	row := db.QueryRow(ctx, `SELECT COALESCE(password_active, false) FROM passwords WHERE guid_id = $1`, guidID)
	err := row.Scan(&active)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return active, nil
}

// CreatePassword создаёт пароль для guid_id
// Принимает уже захешированный пароль (хеш argon2id в виде hex-строки)
// Устанавливает password_active = true
// Обновляет только если пароль ещё неактивен (password_active = false)
func CreatePassword(ctx context.Context, db *pgxpool.Pool, guidID int64, passwordHash string) error {
	// Сначала проверяем, есть ли запись
	var exists bool
	err := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM passwords WHERE guid_id = $1)`, guidID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// Обновляем существующую запись
		_, err = db.Exec(ctx,
			`UPDATE passwords SET 
			   password_hash = $2,
			   password_date = CURRENT_TIMESTAMP,
			   password_active = true
			 WHERE guid_id = $1`,
			guidID, passwordHash)
	} else {
		// Создаем новую запись
		_, err = db.Exec(ctx,
			`INSERT INTO passwords (guid_id, password_hash, password_date, password_active)
			 VALUES ($1, $2, CURRENT_TIMESTAMP, true)`,
			guidID, passwordHash)
	}
	return err
}

// GetPasswordHashByGuid retrieves the password hash for a given guid_id
// Returns hash and found flag
func GetPasswordHashByGuid(ctx context.Context, db *pgxpool.Pool, guidID int64) (string, bool, error) {
	var hash string
	row := db.QueryRow(ctx, "SELECT password_hash FROM passwords WHERE guid_id = $1 LIMIT 1", guidID)
	err := row.Scan(&hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return hash, true, nil
}

// CreateProtection creates or updates protection status for guid_id
func CreateProtection(ctx context.Context, db *pgxpool.Pool, guidID int64, protectionStatus bool) error {
	_, err := db.Exec(ctx,
		`INSERT INTO protection (guid_id, protection_status, protection_date)
		 VALUES ($1, $2, CURRENT_TIMESTAMP)
		 ON CONFLICT (guid_id) DO UPDATE SET 
		   protection_status = $2,
		   protection_date = CURRENT_TIMESTAMP`,
		guidID, protectionStatus)
	return err
}