-- 刪除觸發器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_restaurants_updated_at ON restaurants;
DROP TRIGGER IF EXISTS update_advertisements_updated_at ON advertisements;

-- 刪除觸發器函數
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 刪除索引
DROP INDEX IF EXISTS idx_restaurants_location;
DROP INDEX IF EXISTS idx_restaurants_cuisine;
DROP INDEX IF EXISTS idx_restaurants_rating;
DROP INDEX IF EXISTS idx_favorite_restaurants_user_id;
DROP INDEX IF EXISTS idx_game_sessions_user_id;
DROP INDEX IF EXISTS idx_game_sessions_status;
DROP INDEX IF EXISTS idx_advertisements_active;
DROP INDEX IF EXISTS idx_advertisements_dates;
DROP INDEX IF EXISTS idx_ad_views_ad_id;
DROP INDEX IF EXISTS idx_ad_clicks_ad_id;

-- 刪除資料表（注意順序，先刪除有外鍵的表）
DROP TABLE IF EXISTS ad_clicks;
DROP TABLE IF EXISTS ad_views;
DROP TABLE IF EXISTS game_sessions;
DROP TABLE IF EXISTS advertisements;
DROP TABLE IF EXISTS favorite_restaurants;
DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS user_locations;
DROP TABLE IF EXISTS users;