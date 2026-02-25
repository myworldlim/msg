-- 016_create_secrets.sql
-- Create secrets table for secret words/phrases

CREATE TABLE secrets (
    secret_id BIGSERIAL PRIMARY KEY,
    guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
    secret_hash TEXT NOT NULL,
    secret_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_secrets_guid_id ON secrets(guid_id);
CREATE INDEX idx_secrets_active ON secrets(secret_active);