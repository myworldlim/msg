-- Add verification token columns to passwords table
-- This ensures only users who control their email/phone can register/reset passwords

ALTER TABLE passwords 
ADD COLUMN password_verification_token VARCHAR(255),
ADD COLUMN password_verification_expires TIMESTAMP;

-- Create index for faster token lookups
CREATE INDEX idx_passwords_verification_token ON passwords(password_verification_token) 
WHERE password_verification_token IS NOT NULL;
