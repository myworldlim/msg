package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FindUserUIDByEmail searches users_id for given email
func FindUserUIDByEmail(ctx context.Context, db *pgxpool.Pool, email string) (string, bool, error) {
	var uid string
	row := db.QueryRow(ctx, "SELECT user_uid FROM users_id WHERE user_email = $1 LIMIT 1", email)
	err := row.Scan(&uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return uid, true, nil
}

// FindUserUIDByNumber searches users_id for given phone number
func FindUserUIDByNumber(ctx context.Context, db *pgxpool.Pool, number string) (string, bool, error) {
	var uid string
	row := db.QueryRow(ctx, "SELECT user_uid FROM users_id WHERE user_number = $1 LIMIT 1", number)
	err := row.Scan(&uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return uid, true, nil
}

// FindGuidIDByUserUID returns guid_id and found flag
func FindGuidIDByUserUID(ctx context.Context, db *pgxpool.Pool, userUID string) (int64, bool, error) {
	var guidID int64
	row := db.QueryRow(ctx, "SELECT guid_id FROM guid WHERE user_uid = $1 LIMIT 1", userUID)
	err := row.Scan(&guidID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}
	return guidID, true, nil
}

// GetBlockedStatusByGuidID returns blocked_status (if no record -> false)
func GetBlockedStatusByGuidID(ctx context.Context, db *pgxpool.Pool, guidID int64) (bool, error) {
	var blocked bool
	row := db.QueryRow(ctx, "SELECT blocked_status FROM blocked WHERE guid_id = $1 LIMIT 1", guidID)
	err := row.Scan(&blocked)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return blocked, nil
}

// CreateOrGetUser creates a new user (or returns existing one) in users_id and guid tables.
// email and number can be empty. Returns guid_id and error.
// Uses UNIQUE constraints to handle race conditions safely (postgres will reject duplicates).
func CreateOrGetUser(ctx context.Context, db *pgxpool.Pool, userUID string, email string, number string) (int64, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Convert empty strings to nil for nullable fields
	var emailPtr *string
	var numberPtr *string
	if email != "" {
		emailPtr = &email
	}
	if number != "" {
		numberPtr = &number
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO users_id (user_uid, user_email, user_number, user_reg) VALUES ($1, $2, $3, CURRENT_TIMESTAMP) ON CONFLICT (user_uid) DO NOTHING",
		userUID, emailPtr, numberPtr)
	if err != nil {
		return 0, err
	}

	// Insert into guid with next public ID
	_, err = tx.Exec(ctx,
		"INSERT INTO guid (guid_id, user_uid) VALUES (nextval('guid_public_id_seq'), $1) ON CONFLICT (user_uid) DO NOTHING",
		userUID)
	if err != nil {
		return 0, err
	}

	// Retrieve guid_id
	var guidID int64
	row := tx.QueryRow(ctx, "SELECT guid_id FROM guid WHERE user_uid = $1", userUID)
	err = row.Scan(&guidID)
	if err != nil {
		return 0, err
	}

	// Create password record with password_active = false (user hasn't registered password yet)
	_, err = tx.Exec(ctx,
		"INSERT INTO passwords (guid_id, password_active, password_date) VALUES ($1, false, CURRENT_TIMESTAMP) ON CONFLICT (guid_id) DO NOTHING",
		guidID)
	if err != nil {
		return 0, err
	}

	// Create protection record with protection_status = false (user hasn't set protection yet)
	_, err = tx.Exec(ctx,
		"INSERT INTO protection (guid_id, protection_status, protection_date) VALUES ($1, false, CURRENT_TIMESTAMP) ON CONFLICT (guid_id) DO NOTHING",
		guidID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return guidID, nil
}
