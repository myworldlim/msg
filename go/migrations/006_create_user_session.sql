CREATE TABLE user_session (
    id BIGSERIAL PRIMARY KEY,
    user_uid VARCHAR(36) NOT NULL,
    session_id VARCHAR(100) NOT NULL,
    ip_address VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uid) REFERENCES users(uid)
);