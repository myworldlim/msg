package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CheckPasswordError проверяет статус блокировки пароля
func CheckPasswordError(ctx context.Context, db *pgxpool.Pool, userUid string) (errorActive bool, failedAttempts int, lockedUntil *time.Time, err error) {
	query := `
		SELECT error_active, failed_attempts, locked_until 
		FROM error_password 
		WHERE user_uid = $1
	`
	
	var lockedUntilPtr *time.Time
	err = db.QueryRow(ctx, query, userUid).Scan(&errorActive, &failedAttempts, &lockedUntilPtr)
	
	if err != nil {
		// Если записи нет - пользователь не заблокирован
		if err.Error() == "no rows in result set" {
			return false, 0, nil, nil
		}
		return false, 0, nil, err
	}
	
	// Проверяем не истекла ли блокировка
	if lockedUntilPtr != nil && time.Now().After(*lockedUntilPtr) {
		// Блокировка истекла - сбрасываем статус
		_, updateErr := db.Exec(ctx, `
			UPDATE error_password 
			SET error_active = false, failed_attempts = 0, locked_until = NULL 
			WHERE user_uid = $1
		`, userUid)
		
		if updateErr != nil {
			return false, 0, nil, updateErr
		}
		
		return false, 0, nil, nil
	}
	
	return errorActive, failedAttempts, lockedUntilPtr, nil
}

// CheckPasswordRecovery проверяет доступность восстановления пароля
func CheckPasswordRecovery(ctx context.Context, db *pgxpool.Pool, userUid string) (available bool, method string, contact string, err error) {
	query := `
		SELECT recovery_available, recovery_method, recovery_contact 
		FROM recover_password 
		WHERE user_uid = $1
	`
	
	err = db.QueryRow(ctx, query, userUid).Scan(&available, &method, &contact)
	
	if err != nil {
		// Если записи нет - восстановление недоступно
		if err.Error() == "no rows in result set" {
			return false, "", "", nil
		}
		return false, "", "", err
	}
	
	return available, method, contact, nil
}

// IncrementPasswordError увеличивает счетчик неудачных попыток
func IncrementPasswordError(ctx context.Context, db *pgxpool.Pool, userUid string) error {
	// Получаем текущее количество попыток
	var currentAttempts int
	err := db.QueryRow(ctx, `
		SELECT COALESCE(failed_attempts, 0) 
		FROM error_password 
		WHERE user_uid = $1
	`, userUid).Scan(&currentAttempts)
	
	newAttempts := currentAttempts + 1
	var lockedUntil *time.Time
	
	// Определяем время блокировки по прогрессивной схеме
	if newAttempts >= 3 {
		var lockDuration time.Duration
		switch {
		case newAttempts == 3:
			lockDuration = 5 * time.Minute
		case newAttempts == 4:
			lockDuration = 15 * time.Minute
		case newAttempts == 5:
			lockDuration = 1 * time.Hour
		default: // 6+
			lockDuration = 24 * time.Hour
		}
		
		lockTime := time.Now().Add(lockDuration)
		lockedUntil = &lockTime
	}
	
	// Обновляем или создаем запись
	if err != nil && err.Error() == "no rows in result set" {
		// Создаем новую запись
		_, err = db.Exec(ctx, `
			INSERT INTO error_password (user_uid, failed_attempts, error_active, locked_until, last_attempt)
			VALUES ($1, $2, $3, $4, $5)
		`, userUid, newAttempts, newAttempts >= 3, lockedUntil, time.Now())
	} else {
		// Обновляем существующую запись
		_, err = db.Exec(ctx, `
			UPDATE error_password 
			SET failed_attempts = $2, error_active = $3, locked_until = $4, last_attempt = $5
			WHERE user_uid = $1
		`, userUid, newAttempts, newAttempts >= 3, lockedUntil, time.Now())
	}
	
	return err
}

// ResetPasswordError сбрасывает счетчик при успешном входе
func ResetPasswordError(ctx context.Context, db *pgxpool.Pool, userUid string) error {
	_, err := db.Exec(ctx, `
		UPDATE error_password 
		SET failed_attempts = 0, error_active = false, locked_until = NULL 
		WHERE user_uid = $1
	`, userUid)
	
	return err
}