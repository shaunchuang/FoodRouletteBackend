package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config 包含應用程式的所有配置
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Logger        LoggerConfig
	Auth          AuthConfig
	GoogleAPI     GoogleAPIConfig
	AppleOAuth    AppleOAuthConfig
	Redis         RedisConfig
	Upload        UploadConfig
	CORS          CORSConfig
	RateLimit     RateLimitConfig
	Game          GameConfig
	Advertisement AdvertisementConfig
}

// ServerConfig HTTP 伺服器配置
type ServerConfig struct {
	Port string
	Host string
	Mode string // gin mode: debug, test, release
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoggerConfig 日誌配置
type LoggerConfig struct {
	Level string
}

// AuthConfig JWT 認證配置
type AuthConfig struct {
	Secret    string
	ExpiresIn string
}

// GoogleAPIConfig Google API 配置
type GoogleAPIConfig struct {
	PlacesAPIKey string
	ClientID     string
	ClientSecret string
}

// AppleOAuthConfig Apple OAuth 配置
type AppleOAuthConfig struct {
	ClientID   string
	TeamID     string
	KeyID      string
	PrivateKey string
}

// RedisConfig Redis 配置
type RedisConfig struct {
	URL      string
	Password string
}

// UploadConfig 檔案上傳配置
type UploadConfig struct {
	MaxSize      int64
	AllowedTypes []string
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

// GameConfig 遊戲配置
type GameConfig struct {
	MaxRestaurantsPerRound int
	SessionTimeoutMinutes  int
}

// AdvertisementConfig 廣告配置
type AdvertisementConfig struct {
	ViewCooldownSeconds  int
	ClickCooldownSeconds int
}

// Load 載入配置，優先從環境變數讀取，其次從 .env 檔案
func Load() (*Config, error) {
	// 嘗試載入 .env 檔案（如果存在的話）
	_ = godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		dbPort = 5432
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "food_roulette"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			Secret:    getEnv("JWT_SECRET", "default_secret_change_this"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		GoogleAPI: GoogleAPIConfig{
			PlacesAPIKey: getEnv("GOOGLE_PLACES_API_KEY", ""),
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		},
		AppleOAuth: AppleOAuthConfig{
			ClientID:   getEnv("APPLE_CLIENT_ID", ""),
			TeamID:     getEnv("APPLE_TEAM_ID", ""),
			KeyID:      getEnv("APPLE_KEY_ID", ""),
			PrivateKey: getEnv("APPLE_PRIVATE_KEY", ""),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Upload: UploadConfig{
			MaxSize:      getEnvInt64("UPLOAD_MAX_SIZE", 10485760), // 10MB
			AllowedTypes: strings.Split(getEnv("UPLOAD_ALLOWED_TYPES", "image/jpeg,image/png,image/gif,image/webp"), ","),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
			AllowedMethods: strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
			AllowedHeaders: strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization"), ","),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			Burst:             getEnvInt("RATE_LIMIT_BURST", 10),
		},
		Game: GameConfig{
			MaxRestaurantsPerRound: getEnvInt("GAME_MAX_RESTAURANTS_PER_ROUND", 10),
			SessionTimeoutMinutes:  getEnvInt("GAME_SESSION_TIMEOUT_MINUTES", 30),
		},
		Advertisement: AdvertisementConfig{
			ViewCooldownSeconds:  getEnvInt("AD_VIEW_COOLDOWN_SECONDS", 30),
			ClickCooldownSeconds: getEnvInt("AD_CLICK_COOLDOWN_SECONDS", 60),
		},
	}

	return config, nil
}

// getEnv 取得環境變數，如果不存在則使用預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 取得整數型環境變數，如果不存在或無效則使用預設值
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvInt64 取得 int64 型環境變數，如果不存在或無效則使用預設值
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if int64Value, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int64Value
		}
	}
	return defaultValue
}

// GetDSN 取得資料庫連接字串
func (c *DatabaseConfig) GetDSN() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}
