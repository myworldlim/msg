-- Миграция для создания таблиц безопасности паролей
-- 017_create_password_security.sql

-- Сначала очищаем таблицу passwords от ненужных полей
ALTER TABLE passwords 
DROP COLUMN IF EXISTS password_reset,
DROP COLUMN IF EXISTS password_attempts,
DROP COLUMN IF EXISTS password_verification,
DROP COLUMN IF EXISTS password_last_date,
DROP COLUMN IF EXISTS password_verification_token,
DROP COLUMN IF EXISTS password_verification_expires;

-- Создаем таблицу для отслеживания ошибок ввода пароля
CREATE TABLE IF NOT EXISTS error_password (
  error_id BIGSERIAL PRIMARY KEY,
  guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
  user_uid VARCHAR(36) NOT NULL,
  failed_attempts INT DEFAULT 0,
  locked_until TIMESTAMPTZ,
  last_attempt TIMESTAMPTZ,
  error_active BOOLEAN DEFAULT false, -- false при регистрации, true при первой ошибке
  created_at TIMESTAMPTZ DEFAULT now(),
  FOREIGN KEY (user_uid) REFERENCES users_id(user_uid) ON DELETE CASCADE
);

-- Индексы для error_password
CREATE UNIQUE INDEX IF NOT EXISTS idx_error_password_guid ON error_password (guid_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_error_password_useruid ON error_password (user_uid);

-- Создаем таблицу для восстановления паролей
CREATE TABLE IF NOT EXISTS recover_password (
  recover_id BIGSERIAL PRIMARY KEY,
  guid_id BIGINT NOT NULL REFERENCES guid(guid_id) ON DELETE CASCADE,
  user_uid VARCHAR(36) NOT NULL,
  recovery_available BOOLEAN DEFAULT false, -- false при регистрации
  recovery_method VARCHAR(10), -- 'email', 'sms', null
  recovery_contact VARCHAR(255), -- замаскированный контакт для показа
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  FOREIGN KEY (user_uid) REFERENCES users_id(user_uid) ON DELETE CASCADE
);

-- Индексы для recover_password
CREATE UNIQUE INDEX IF NOT EXISTS idx_recover_password_guid ON recover_password (guid_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_recover_password_useruid ON recover_password (user_uid);

-- Комментарии к таблицам
COMMENT ON TABLE error_password IS 'Отслеживание неудачных попыток ввода пароля и блокировок';
COMMENT ON TABLE recover_password IS 'Настройки восстановления паролей для пользователей';

-- Комментарии к полям
COMMENT ON COLUMN error_password.error_active IS 'Активна ли блокировка (false при регистрации)';
COMMENT ON COLUMN error_password.failed_attempts IS 'Количество неудачных попыток ввода пароля';
COMMENT ON COLUMN error_password.locked_until IS 'Время до которого заблокирован аккаунт';

COMMENT ON COLUMN recover_password.recovery_available IS 'Доступно ли восстановление пароля';
COMMENT ON COLUMN recover_password.recovery_method IS 'Метод восстановления: email или sms';
COMMENT ON COLUMN recover_password.recovery_contact IS 'Замаскированный контакт для показа пользователю';







-- Заполняем таблицу error_password (все пользователи начинают с error_active = false)
INSERT INTO error_password (guid_id, user_uid, failed_attempts, error_active) VALUES
(1000000, 'ZE5vaJpZ3SJFPzG7SuRv', 0, false),
(1000001, 'VLIQLB8OlvP57BGikiR9', 0, false),
(1000002, 'jDUFptJCxYpa70rZhlD7', 0, false),
(1000003, 'LMlOc5sDSxE4KIORPCeX', 0, false),
(1000004, 'MPv7mjzhksIDnGteMayq', 0, false),
(1000005, 'E9C9FsIxp9fx7kG3A344', 0, false),
(1000006, 'bfEe02Ut9WM39zltkKuT', 0, false),
(1000007, 'wecKBVsoqud6DTS9G9V3', 0, false),
(1000008, 'NX3iIaDdKfo3uxCuBV0Q', 0, false),
(1000009, '3s9q0oDPg8Ovi2glng9E', 0, false),
(1000010, 'P6xC1nzsPSJ4q3vnVgoK', 0, false),
(1000011, 'yWwL91ZTGLWyWrnLANMB', 0, false),
(1000012, 'POhciNIpvWZM6ryPhMAm', 0, false),
(1000013, 'upTKmsmWQJgTkS5snFF2', 0, false),
(1000014, 'K34SZxe4tirzQMAbDm85', 0, false),
(1000015, 'mDdhjE3U9mRFyOuA6JiT', 0, false),
(1000016, 'ifw26HqpWm5c0vuPT3u6', 0, false),
(1000017, 'FuUtOp8yaDAWbF2bJh7R', 0, false),
(1000018, 'wJeu6730X9JSGEwWBgMw', 0, false),
(1000019, 'Xy7RZsx8xzrOrGAxPZSO', 0, false),
(1000020, 'T9mi1d1KlqI9OpyFLEL0', 0, false),
(1000021, 'jHXjMsUeSsD45cXlnwjw', 0, false),
(1000022, 'KxlWFEfuN9AmsaupYGUl', 0, false),
(1000023, 'zN1D6BoVNmDJ2io3tX2x', 0, false),
(1000024, 'GPxTZBkXK1kUh76AuCMM', 0, false),
(1000025, 'E04m9n4xkeiG5rjhKDTJ', 0, false);

-- Заполняем таблицу recover_password (все пользователи начинают с recovery_available = false)
INSERT INTO recover_password (guid_id, user_uid, recovery_available) VALUES
(1000000, 'ZE5vaJpZ3SJFPzG7SuRv', false),
(1000001, 'VLIQLB8OlvP57BGikiR9', false),
(1000002, 'jDUFptJCxYpa70rZhlD7', false),
(1000003, 'LMlOc5sDSxE4KIORPCeX', false),
(1000004, 'MPv7mjzhksIDnGteMayq', false),
(1000005, 'E9C9FsIxp9fx7kG3A344', false),
(1000006, 'bfEe02Ut9WM39zltkKuT', false),
(1000007, 'wecKBVsoqud6DTS9G9V3', false),
(1000008, 'NX3iIaDdKfo3uxCuBV0Q', false),
(1000009, '3s9q0oDPg8Ovi2glng9E', false),
(1000010, 'P6xC1nzsPSJ4q3vnVgoK', false),
(1000011, 'yWwL91ZTGLWyWrnLANMB', false),
(1000012, 'POhciNIpvWZM6ryPhMAm', false),
(1000013, 'upTKmsmWQJgTkS5snFF2', false),
(1000014, 'K34SZxe4tirzQMAbDm85', false),
(1000015, 'mDdhjE3U9mRFyOuA6JiT', false),
(1000016, 'ifw26HqpWm5c0vuPT3u6', false),
(1000017, 'FuUtOp8yaDAWbF2bJh7R', false),
(1000018, 'wJeu6730X9JSGEwWBgMw', false),
(1000019, 'Xy7RZsx8xzrOrGAxPZSO', false),
(1000020, 'T9mi1d1KlqI9OpyFLEL0', false),
(1000021, 'jHXjMsUeSsD45cXlnwjw', false),
(1000022, 'KxlWFEfuN9AmsaupYGUl', false),
(1000023, 'zN1D6BoVNmDJ2io3tX2x', false),
(1000024, 'GPxTZBkXK1kUh76AuCMM', false),
(1000025, 'E04m9n4xkeiG5rjhKDTJ', false);


-- Обновляем некоторых пользователей с разными статусами для тестирования

-- Пользователь с восстановлением через email
UPDATE recover_password 
SET recovery_available = true, recovery_method = 'email', recovery_contact = 'my***@gmail.com'
WHERE user_uid = 'jDUFptJCxYpa70rZhlD7';

-- Пользователь с восстановлением через SMS
UPDATE recover_password 
SET recovery_available = true, recovery_method = 'sms', recovery_contact = '+7995***5002'
WHERE user_uid = 'LMlOc5sDSxE4KIORPCeX';

-- Заблокированный пользователь на 5 минут (3 попытки)
UPDATE error_password 
SET failed_attempts = 3, error_active = true, locked_until = NOW() + INTERVAL '5 minutes', last_attempt = NOW()
WHERE user_uid = 'E9C9FsIxp9fx7kG3A344';

-- Заблокированный пользователь на 15 минут (4 попытки)
UPDATE error_password 
SET failed_attempts = 4, error_active = true, locked_until = NOW() + INTERVAL '15 minutes', last_attempt = NOW()
WHERE user_uid = 'bfEe02Ut9WM39zltkKuT';

-- Заблокированный пользователь на 1 час (5 попыток)
UPDATE error_password 
SET failed_attempts = 5, error_active = true, locked_until = NOW() + INTERVAL '1 hour', last_attempt = NOW()
WHERE user_uid = 'wecKBVsoqud6DTS9G9V3';

-- Заблокированный пользователь на 24 часа (6+ попыток)
UPDATE error_password 
SET failed_attempts = 7, error_active = true, locked_until = NOW() + INTERVAL '24 hours', last_attempt = NOW()
WHERE user_uid = 'NX3iIaDdKfo3uxCuBV0Q';

-- Пользователь с 2 неудачными попытками (еще не заблокирован)
UPDATE error_password 
SET failed_attempts = 2, error_active = false, last_attempt = NOW()
WHERE user_uid = 'MPv7mjzhksIDnGteMayq';

-- Пользователь с восстановлением + заблокирован
UPDATE recover_password 
SET recovery_available = true, recovery_method = 'email', recovery_contact = 'test***@example.com'
WHERE user_uid = 'wecKBVsoqud6DTS9G9V3';

-- Проверяем результат
SELECT 'error_password' as table_name, user_uid, failed_attempts, error_active, 
       CASE 
         WHEN locked_until IS NULL THEN 'не заблокирован'
         WHEN locked_until > NOW() THEN 'заблокирован до ' || locked_until::text
         ELSE 'блокировка истекла'
       END as lock_status
FROM error_password 
WHERE failed_attempts > 0 OR error_active = true
ORDER BY failed_attempts DESC;

SELECT 'recover_password' as table_name, user_uid, recovery_available, recovery_method, recovery_contact
FROM recover_password 
WHERE recovery_available = true;

