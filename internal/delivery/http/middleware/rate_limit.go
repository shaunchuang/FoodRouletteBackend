package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware 速率限制中介軟體
func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "請求過於頻繁，請稍後再試",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// PerIPRateLimitMiddleware 基於 IP 的速率限制中介軟體
func PerIPRateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
			limiters[ip] = limiter
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "請求過於頻繁，請稍後再試",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
