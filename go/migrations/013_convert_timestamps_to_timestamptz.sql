-- 013_convert_timestamps_to_timestamptz.sql
-- Convert some TIMESTAMP (without time zone) columns to TIMESTAMPTZ
-- We assume existing values were recorded in Europe/Moscow (GMT+3) local time and
-- want to preserve the same wall-clock time but store it with +03 offset (timestamptz).
-- IMPORTANT: make a DB backup before running this migration.

BEGIN;

-- users_id.user_reg: convert to timestamptz, interpreting existing values as Europe/Moscow local time
ALTER TABLE users_id
  ALTER COLUMN user_reg TYPE TIMESTAMPTZ
  USING user_reg AT TIME ZONE 'Europe/Moscow';
ALTER TABLE users_id
  ALTER COLUMN user_reg SET DEFAULT now();

-- blocked.blocked_date
ALTER TABLE blocked
  ALTER COLUMN blocked_date TYPE TIMESTAMPTZ
  USING blocked_date AT TIME ZONE 'Europe/Moscow';
ALTER TABLE blocked
  ALTER COLUMN blocked_date SET DEFAULT now();

-- password_verification_expires (added previously as TIMESTAMP) -> convert to timestamptz
ALTER TABLE passwords
  ALTER COLUMN password_verification_expires TYPE TIMESTAMPTZ
  USING password_verification_expires AT TIME ZONE 'Europe/Moscow';

COMMIT;

-- After running migration, you can verify with:
-- SELECT id, user_uid, user_reg FROM users_id ORDER BY id DESC LIMIT 5;
-- SELECT guid_id, password_date, password_last_date, password_verification_expires FROM passwords ORDER BY password_id DESC LIMIT 5;
