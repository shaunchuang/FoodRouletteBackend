-- 移除權限相關表
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;

-- 移除使用者表的新增欄位
ALTER TABLE users 
DROP COLUMN IF EXISTS role,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS email_verified,
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS failed_login_attempts,
DROP COLUMN IF EXISTS locked_until;
