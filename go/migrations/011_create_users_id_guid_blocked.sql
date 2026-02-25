-- 011_create_users_id_guid_blocked.sql

CREATE TABLE IF NOT EXISTS users_id (
    id BIGSERIAL PRIMARY KEY,
    user_uid VARCHAR(36) UNIQUE NOT NULL,
    user_number VARCHAR(32) UNIQUE,
    user_email VARCHAR(255) UNIQUE,
    user_reg TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_id_email ON users_id (user_email);
CREATE INDEX IF NOT EXISTS idx_users_id_number ON users_id (user_number);

CREATE TABLE IF NOT EXISTS guid (
    guid_id BIGSERIAL PRIMARY KEY,
    user_uid VARCHAR(36) UNIQUE NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_guid_useruid ON guid (user_uid);

CREATE TABLE IF NOT EXISTS blocked (
    blocked_id BIGSERIAL PRIMARY KEY,
    guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
    blocked_status BOOLEAN NOT NULL DEFAULT false,
    blocked_type VARCHAR(64),
    blocked_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_blocked_guid_date ON blocked (guid_id, blocked_date DESC);



CREATE TABLE passwords (
  password_id BIGSERIAL PRIMARY KEY,
  guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
  password_hash VARCHAR(64),         -- хеш пароля (например, hex-строка SHA256/argon2id/bcrypt)
  password_reset VARCHAR(64),        -- хеш/токен для сброса пароля (опционально)
  password_attempts INT DEFAULT 0,   -- количество попыток ввода пароля (для блокировки)
  password_verification INT DEFAULT 0, -- количество успешных верификаций (можно для статистики)
  password_date TIMESTAMPTZ NOT NULL DEFAULT now(), -- дата создания пароля
  password_last_date TIMESTAMPTZ,    -- дата последнего изменения пароля
  password_active BOOLEAN NOT NULL DEFAULT false    -- false при создании пользователя, true после успешной регистрации пароля
);
CREATE UNIQUE INDEX idx_passwords_guid ON passwords (guid_id);
-- Add verification token columns to passwords table
-- This ensures only users who control their email/phone can register/reset passwords

ALTER TABLE passwords 
ADD COLUMN password_verification_token VARCHAR(255),
ADD COLUMN password_verification_expires TIMESTAMP;

-- Create index for faster token lookups
CREATE INDEX idx_passwords_verification_token ON passwords(password_verification_token) 
WHERE password_verification_token IS NOT NULL;