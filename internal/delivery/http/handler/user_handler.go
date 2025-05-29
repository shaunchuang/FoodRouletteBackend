package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/internal/usecase"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// UserHandler 使用者 HTTP 處理器
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler 建立使用者處理器
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// Register 使用者註冊
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("註冊請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "註冊成功",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// Login 使用者登入
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("登入請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	token, err := h.userUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"token":   token,
	})
}

// GetProfile 取得使用者資料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	user, err := h.userUseCase.GetProfile(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Username,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// UpdateLocation 更新使用者位置
func (h *UserHandler) UpdateLocation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	var req domain.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新位置請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	if err := h.userUseCase.UpdateLocation(c.Request.Context(), userID.(int), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "位置更新成功",
	})
}

// GetLocation 取得使用者位置
func (h *UserHandler) GetLocation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	location, err := h.userUseCase.GetLocation(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"location": location,
	})
}