CREATE TABLE user_password (
    password_id BIGSERIAL PRIMARY KEY,
    user_uid VARCHAR(36) NOT NULL UNIQUE,
    password VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uid) REFERENCES users(uid)
);