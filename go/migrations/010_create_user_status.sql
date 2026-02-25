CREATE TABLE user_status_ws (
    id BIGSERIAL PRIMARY KEY,
    user_uid VARCHAR(36) UNIQUE NOT NULL,
    connection_id VARCHAR(36) UNIQUE NOT NULL,
    status VARCHAR(10) CHECK (status IN ('online', 'offline', 'idle')) DEFAULT 'offline',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_ping_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_uid) REFERENCES users(uid)
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_status_ws_updated_at
    BEFORE UPDATE ON user_status_ws
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();