package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/internal/usecase"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// AuthHandler 認證處理器
type AuthHandler struct {
	userUseCase usecase.UserService
	authService usecase.AuthService
}

// NewAuthHandler 建立認證處理器
func NewAuthHandler(userUseCase usecase.UserService, authService usecase.AuthService) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUseCase,
		authService: authService,
	}
}

// GoogleLogin Google 登入
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req domain.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Google 登入請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "請求參數錯誤",
			},
		})
		return
	}

	// 驗證 Google ID Token
	userInfo, err := h.authService.VerifyGoogleToken(req.IDToken)
	if err != nil {
		logger.Error("Google Token 驗證失敗", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "無效的 Google ID Token",
			},
		})
		return
	}

	// 處理使用者登入/註冊
	user, token, err := h.userUseCase.OAuthLogin(c.Request.Context(), userInfo)
	if err != nil {
		logger.Error("Google OAuth 登入失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "LOGIN_FAILED",
				"message": "登入失敗",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":             user.ID,
				"email":          user.Email,
				"username":       user.Username,
				"avatar":         user.Avatar,
				"role":           user.Role,
				"status":         user.Status,
				"provider":       user.Provider,
				"email_verified": user.EmailVerified,
			},
		},
	})
}

// AppleLogin Apple ID 登入
func (h *AuthHandler) AppleLogin(c *gin.Context) {
	var req domain.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Apple 登入請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "請求參數錯誤",
			},
		})
		return
	}

	// 驗證 Apple ID Token
	userInfo, err := h.authService.VerifyAppleToken(req.IDToken)
	if err != nil {
		logger.Error("Apple Token 驗證失敗", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "無效的 Apple ID Token",
			},
		})
		return
	}

	// 處理使用者登入/註冊
	user, token, err := h.userUseCase.OAuthLogin(c.Request.Context(), userInfo)
	if err != nil {
		logger.Error("Apple OAuth 登入失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "LOGIN_FAILED",
				"message": "登入失敗",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":             user.ID,
				"email":          user.Email,
				"username":       user.Username,
				"avatar":         user.Avatar,
				"role":           user.Role,
				"status":         user.Status,
				"provider":       user.Provider,
				"email_verified": user.EmailVerified,
			},
		},
	})
}
