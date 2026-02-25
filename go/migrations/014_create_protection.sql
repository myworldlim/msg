-- 014_create_protection.sql
-- Create protection table for storing user protection preferences (2FA, secret questions, etc.)

CREATE TABLE protection (
    protection_id BIGSERIAL PRIMARY KEY,
    guid_id BIGINT NOT NULL UNIQUE REFERENCES guid(guid_id) ON DELETE CASCADE,
    protection_status BOOLEAN NOT NULL DEFAULT false,
    protection_date TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_protection_guid ON protection(guid_id);
