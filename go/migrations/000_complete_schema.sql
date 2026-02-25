-- Complete database schema for ChitChat
-- All tables with final corrections applied

-- Users identification table
CREATE TABLE users_id (
    id BIGSERIAL PRIMARY KEY,
    user_uid VARCHAR(36) UNIQUE NOT NULL,
    user_number VARCHAR(32) UNIQUE,
    user_email VARCHAR(255) UNIQUE,
    user_reg TIMESTAMPTZ DEFAULT now()
);

CREATE UNIQUE INDEX idx_users_email ON users_id (LOWER(user_email));
CREATE UNIQUE INDEX idx_users_number ON users_id (user_number);

-- GUID mapping table
CREATE TABLE guid (
    id BIGSERIAL PRIMARY KEY,                    -- просто порядковый номер PK его ни где не учитываем
    guid_id BIGINT UNIQUE NOT NULL,              -- Внутренний ID (начинаем с 1000000)
    user_uid VARCHAR(36) UNIQUE NOT NULL         -- из таблицы user_id поля user_uid
);

-- Устанавливаем начальное значение для guid_id
CREATE SEQUENCE guid_public_id_seq START 1000000;

CREATE UNIQUE INDEX idx_guid_useruid ON guid (user_uid);
CREATE UNIQUE INDEX idx_guid_public_id ON guid (guid_id);


-- User blocking table
CREATE TABLE blocked (
    blocked_id BIGSERIAL PRIMARY KEY,
    guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
    blocked_status BOOLEAN NOT NULL DEFAULT false,
    blocked_type VARCHAR(64),
    blocked_date TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_blocked_guid_date ON blocked (guid_id, blocked_date DESC);

-- Passwords table (corrected hash field size)
CREATE TABLE passwords (
    password_id BIGSERIAL PRIMARY KEY,
    guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
    password_hash VARCHAR(128),                        -- Increased size for salt+hash
    password_reset VARCHAR(64),
    password_attempts INT DEFAULT 0,
    password_verification INT DEFAULT 0,
    password_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    password_last_date TIMESTAMPTZ,
    password_active BOOLEAN NOT NULL DEFAULT false,
    password_verification_token VARCHAR(255),
    password_verification_expires TIMESTAMPTZ,
    UNIQUE(guid_id)
);

CREATE INDEX idx_passwords_verification_token ON passwords(password_verification_token) 
WHERE password_verification_token IS NOT NULL;

-- Protection settings table
CREATE TABLE protection (
    protection_id BIGSERIAL PRIMARY KEY,
    guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
    protection_status BOOLEAN NOT NULL DEFAULT false,
    protection_date TIMESTAMPTZ DEFAULT now(),
    UNIQUE(guid_id)
);

CREATE INDEX idx_protection_guid ON protection(guid_id);

-- Sessions table
CREATE TABLE sessions (
    session_id BIGSERIAL PRIMARY KEY,
    user_uid TEXT NOT NULL,
    guid_id BIGINT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
    session_token TEXT NOT NULL UNIQUE,
    session_refresh TEXT NOT NULL UNIQUE,
    session_user_agent TEXT,
    ip_address TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    session_expires_at TIMESTAMPTZ,
    session_refresh_expires_at TIMESTAMPTZ
);

CREATE INDEX idx_sessions_user_uid ON sessions(user_uid);
CREATE INDEX idx_sessions_token ON sessions(session_token);
CREATE INDEX idx_sessions_refresh ON sessions(session_refresh);

-- Secrets table (with unique constraint for guid_id)
CREATE TABLE secrets (
    secret_id BIGSERIAL PRIMARY KEY,
    guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
    secret_hash TEXT NOT NULL,
    secret_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(guid_id)
);

CREATE INDEX idx_secrets_guid_id ON secrets(guid_id);
CREATE INDEX idx_secrets_active ON secrets(secret_active);