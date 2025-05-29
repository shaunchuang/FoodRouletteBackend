package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/internal/usecase"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// GameHandler 遊戲 HTTP 處理器
type GameHandler struct {
	gameUseCase *usecase.GameUseCase
}

// NewGameHandler 建立遊戲處理器
func NewGameHandler(gameUseCase *usecase.GameUseCase) *GameHandler {
	return &GameHandler{
		gameUseCase: gameUseCase,
	}
}

// StartGame 開始遊戲
func (h *GameHandler) StartGame(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	var req domain.StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("開始遊戲請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	// 設定預設值
	if req.Radius == 0 {
		req.Radius = 1000 // 預設 1 公里
	}

	session, err := h.gameUseCase.StartGame(c.Request.Context(), userID.(int), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "遊戲開始成功",
		"session": session,
	})
}

// CompleteGame 完成遊戲
func (h *GameHandler) CompleteGame(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	var req domain.CompleteGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("完成遊戲請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	result, err := h.gameUseCase.CompleteGame(c.Request.Context(), userID.(int), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "遊戲完成",
		"result":  result,
	})
}

// GetGameHistory 取得遊戲歷史
func (h *GameHandler) GetGameHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	sessions, err := h.gameUseCase.GetGameHistory(c.Request.Context(), userID.(int), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}