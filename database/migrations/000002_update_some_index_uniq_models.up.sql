-- =============================================
-- departments.code
-- Tên có thể là: departments_code_key (PostgreSQL) hoặc uni_departments_code (GORM)
-- =============================================
ALTER TABLE departments DROP CONSTRAINT IF EXISTS departments_code_key;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS uni_departments_code;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS idx_departments_code;
DROP INDEX IF EXISTS idx_departments_code;

CREATE UNIQUE INDEX IF NOT EXISTS idx_departments_code_active
    ON departments (code)
    WHERE deleted_at IS NULL;

-- =============================================
-- services.code
-- =============================================
ALTER TABLE services DROP CONSTRAINT IF EXISTS services_code_key;
ALTER TABLE services DROP CONSTRAINT IF EXISTS uni_services_code;
ALTER TABLE services DROP CONSTRAINT IF EXISTS idx_services_code;
DROP INDEX IF EXISTS idx_services_code;

CREATE UNIQUE INDEX IF NOT EXISTS idx_services_code_active
    ON services (code)
    WHERE deleted_at IS NULL;

-- =============================================
-- users.citizen_id
-- GORM uniqueIndex tag tạo constraint có cùng tên với index
-- Phải DROP CONSTRAINT thay vì DROP INDEX
-- =============================================
ALTER TABLE users DROP CONSTRAINT IF EXISTS idx_users_citizen_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_citizen_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_citizen_id_key;
DROP INDEX IF EXISTS idx_users_citizen_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_citizen_id_active
    ON users (citizen_id)
    WHERE deleted_at IS NULL;

-- =============================================
-- users.email
-- =============================================
ALTER TABLE users DROP CONSTRAINT IF EXISTS idx_users_email;
ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_email;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
DROP INDEX IF EXISTS idx_users_email;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active
    ON users (email)
    WHERE deleted_at IS NULL;
