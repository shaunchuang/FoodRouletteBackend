package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/usecase"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// AuthService 認證服務介面
type AuthService interface {
	ValidateToken(token string) (int, error)
}

// AuthMiddleware 認證中介軟體
func AuthMiddleware(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Header 取得 Authorization token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少認證 token",
			})
			c.Abort()
			return
		}

		// 檢查 Bearer token 格式
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "無效的 token 格式",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// 驗證 JWT token
		userID, err := authService.ValidateToken(token)
		if err != nil {
			logger.Error("token 驗證失敗", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "無效的 token",
			})
			c.Abort()
			return
		}

		// 將使用者 ID 存入 context
		c.Set("user_id", userID)

		logger.Debug("使用者認證成功", zap.Int("user_id", userID))
		c.Next()
	}
}

// AdminMiddleware 管理員權限中介軟體
func AdminMiddleware(userService usecase.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未認證的使用者",
			})
			c.Abort()
			return
		}

		// 從資料庫獲取使用者資訊
		user, err := userService.GetProfile(c.Request.Context(), userID.(int))
		if err != nil {
			logger.Error("取得使用者資料失敗", zap.Error(err), zap.Any("user_id", userID))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "取得使用者資料失敗",
			})
			c.Abort()
			return
		}

		// 檢查使用者狀態
		if !user.IsActive() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "使用者帳號已被停用",
			})
			c.Abort()
			return
		}

		// 檢查是否被鎖定
		if user.IsLocked() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "使用者帳號已被鎖定",
			})
			c.Abort()
			return
		}

		// 檢查管理員權限
		if !user.HasAdminAccess() {
			logger.Warn("非管理員使用者嘗試存取管理功能",
				zap.Int("user_id", userID.(int)),
				zap.String("role", string(user.Role)))
			c.JSON(http.StatusForbidden, gin.H{
				"error": "權限不足，需要管理員權限",
			})
			c.Abort()
			return
		}

		// 將使用者資訊存入 context
		c.Set("user", user)
		logger.Debug("管理員權限驗證成功",
			zap.Int("user_id", userID.(int)),
			zap.String("role", string(user.Role)))
		c.Next()
	}
}

// OptionalAuthMiddleware 可選認證中介軟體 (for anonymous users)
func OptionalAuthMiddleware(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 沒有 token，繼續執行但不設置使用者 ID
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// token 格式錯誤，但不阻止請求
			c.Next()
			return
		}

		token := tokenParts[1]
		userID, err := authService.ValidateToken(token)
		if err != nil {
			// token 無效，但不阻止請求
			logger.Debug("可選認證失敗", zap.Error(err))
			c.Next()
			return
		}

		// 成功驗證則設置使用者 ID
		c.Set("user_id", userID)
		logger.Debug("可選認證成功", zap.Int("user_id", userID))
		c.Next()
	}
}
