package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetSecretHashByUserUID получает хеш секретного слова по user_uid
func GetSecretHashByUserUID(ctx context.Context, db *pgxpool.Pool, userUID string) (string, bool, error) {
	var hash string
	row := db.QueryRow(ctx, 
		`SELECT s.secret_hash 
		 FROM secrets s 
		 JOIN guid g ON s.guid_id = g.guid_id 
		 WHERE g.user_uid = $1 AND s.secret_active = true 
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