package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// AuthMiddleware 認證中介軟體
func AuthMiddleware() gin.HandlerFunc {
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

		// TODO: 實作 JWT token 驗證
		// 這裡應該驗證 JWT token 並取得使用者 ID
		// userID, err := authService.ValidateToken(token)
		// if err != nil {
		//     c.JSON(http.StatusUnauthorized, gin.H{
		//         "error": "無效的 token",
		//     })
		//     c.Abort()
		//     return
		// }

		// 暫時的假驗證，實際應用中需要實作 JWT 驗證
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "無效的 token",
			})
			c.Abort()
			return
		}

		// 將使用者 ID 存入 context（暫時設為 1）
		// 實際應用中應該從 JWT token 中解析出真實的使用者 ID
		c.Set("user_id", 1)

		logger.Debug("使用者認證成功", zap.Int("user_id", 1))
		c.Next()
	}
}

// AdminMiddleware 管理員權限中介軟體
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未認證的使用者",
			})
			c.Abort()
			return
		}

		// TODO: 檢查使用者是否為管理員
		// 這裡應該查詢資料庫確認使用者權限
		// isAdmin, err := userService.IsAdmin(userID)
		// if err != nil || !isAdmin {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "error": "沒有管理員權限",
		//     })
		//     c.Abort()
		//     return
		// }

		logger.Debug("管理員權限驗證成功", zap.Any("user_id", userID))
		c.Next()
	}
}