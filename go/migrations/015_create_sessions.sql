-- 015_create_sessions.sql
-- Create sessions table to store session and refresh tokens

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
