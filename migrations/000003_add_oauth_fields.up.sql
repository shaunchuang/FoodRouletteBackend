-- 新增 OAuth 相關欄位
ALTER TABLE users ADD COLUMN provider VARCHAR(20) DEFAULT 'local';
ALTER TABLE users ADD COLUMN provider_id VARCHAR(255);
ALTER TABLE users ADD COLUMN avatar TEXT;

-- 修改 password 欄位為可選（第三方登入時可為 NULL）
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;

-- 建立索引
CREATE INDEX idx_users_provider_id ON users(provider, provider_id);
CREATE UNIQUE INDEX idx_users_provider_unique ON users(provider, provider_id) WHERE provider IS NOT NULL AND provider != 'local';

-- 更新現有使用者的 provider 為 'local'
UPDATE users SET provider = 'local' WHERE provider IS NULL;
