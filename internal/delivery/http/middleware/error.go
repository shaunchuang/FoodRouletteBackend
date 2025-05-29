package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// ErrorMiddleware 錯誤處理中介軟體
func ErrorMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// 處理錯誤
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			logger.Error("HTTP Request Error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
			)

			// 根據錯誤類型返回適當的 HTTP 狀態碼
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "請求參數錯誤",
					"details": err.Error(),
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "內部伺服器錯誤",
				})
			}
		}
	})
}