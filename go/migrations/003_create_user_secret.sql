CREATE TABLE user_secret (
    secret_id BIGSERIAL PRIMARY KEY,
    user_uid VARCHAR(20) NOT NULL,
    secret VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uid) REFERENCES users(user_uid)
);