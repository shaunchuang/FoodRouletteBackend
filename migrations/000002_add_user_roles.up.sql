-- 為使用者表添加角色和狀態欄位
ALTER TABLE users 
ADD COLUMN role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
ADD COLUMN status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
ADD COLUMN email_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN last_login_at TIMESTAMP,
ADD COLUMN failed_login_attempts INTEGER DEFAULT 0,
ADD COLUMN locked_until TIMESTAMP;

-- 創建使用者角色索引
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_email_verified ON users(email_verified);

-- 創建預設管理員使用者（密碼應該在應用程式中加密）
-- 這裡只是示例，實際部署時應該通過應用程式創建
INSERT INTO users (email, username, password, role, status, email_verified) 
VALUES ('admin@foodroulette.com', 'admin', '$2a$10$example_hashed_password', 'admin', 'active', TRUE)
ON CONFLICT (email) DO NOTHING;

-- 創建權限相關表（可選，用於更細粒度的權限控制）
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 創建角色權限關聯表
CREATE TABLE IF NOT EXISTS role_permissions (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role, permission_id)
);

-- 插入基本權限
INSERT INTO permissions (name, description) VALUES
('user.read', '讀取使用者資料'),
('user.write', '修改使用者資料'),
('user.delete', '刪除使用者'),
('restaurant.read', '讀取餐廳資料'),
('restaurant.write', '修改餐廳資料'),
('restaurant.delete', '刪除餐廳'),
('advertisement.read', '讀取廣告資料'),
('advertisement.write', '修改廣告資料'),
('advertisement.delete', '刪除廣告'),
('game.read', '讀取遊戲資料'),
('game.manage', '管理遊戲'),
('admin.dashboard', '存取管理後台'),
('system.config', '系統配置')
ON CONFLICT (name) DO NOTHING;

-- 為管理員角色分配權限
INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions
ON CONFLICT (role, permission_id) DO NOTHING;

-- 為一般使用者分配基本權限
INSERT INTO role_permissions (role, permission_id)
SELECT 'user', id FROM permissions 
WHERE name IN ('user.read', 'restaurant.read', 'game.read')
ON CONFLICT (role, permission_id) DO NOTHING;

-- 創建索引
CREATE INDEX idx_role_permissions_role ON role_permissions(role);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);
