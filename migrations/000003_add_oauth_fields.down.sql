-- 移除 OAuth 相關欄位
DROP INDEX IF EXISTS idx_users_provider_unique;
DROP INDEX IF EXISTS idx_users_provider_id;

ALTER TABLE users DROP COLUMN IF EXISTS avatar;
ALTER TABLE users DROP COLUMN IF EXISTS provider_id;
ALTER TABLE users DROP COLUMN IF EXISTS provider;

-- 恢復 password 欄位為必要
ALTER TABLE users ALTER COLUMN password SET NOT NULL;
