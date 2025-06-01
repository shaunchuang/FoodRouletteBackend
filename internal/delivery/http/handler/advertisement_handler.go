package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// AdvertisementHandler 廣告 HTTP 處理器
type AdvertisementHandler struct {
	adUseCase AdvertisementUseCase
}

// AdvertisementUseCase 廣告用例介面
type AdvertisementUseCase interface {
	GetActiveAdvertisements(ctx context.Context, limit int) ([]domain.Advertisement, error)
	GetAdvertisementByID(ctx context.Context, id int) (*domain.Advertisement, error)
	RecordAdView(ctx context.Context, view *domain.AdView) error
	RecordAdClick(ctx context.Context, click *domain.AdClick) error
	GetAdvertisementStatistics(ctx context.Context, adID int, period string) (*domain.AdStatistics, error)
	CreateAd(ctx context.Context, req *domain.CreateAdRequest) error
	UpdateAd(ctx context.Context, adID int, req *domain.CreateAdRequest) error
	DeleteAd(ctx context.Context, adID int) error
	GetAllAds(ctx context.Context, limit, offset int) ([]domain.Advertisement, error)
}

// NewAdvertisementHandler 建立廣告處理器
func NewAdvertisementHandler(adUseCase AdvertisementUseCase) *AdvertisementHandler {
	return &AdvertisementHandler{
		adUseCase: adUseCase,
	}
}

// GetActiveAds 取得活躍廣告
func (h *AdvertisementHandler) GetActiveAds(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // 預設值
	}

	ads, err := h.adUseCase.GetActiveAdvertisements(c.Request.Context(), limit)
	if err != nil {
		logger.Error("取得活躍廣告失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "取得廣告失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"advertisements": ads,
		"count":          len(ads),
	})
}

// GetStatistics 取得廣告統計
func (h *AdvertisementHandler) GetStatistics(c *gin.Context) {
	idStr := c.Param("id")
	adID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的廣告 ID",
		})
		return
	}

	period := c.DefaultQuery("period", "daily")

	stats, err := h.adUseCase.GetAdvertisementStatistics(c.Request.Context(), adID, period)
	if err != nil {
		logger.Error("取得廣告統計失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "取得統計失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statistics": stats,
	})
}

// CreateAd 建立廣告（管理功能）
func (h *AdvertisementHandler) CreateAd(c *gin.Context) {
	var req domain.CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("建立廣告請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	if err := h.adUseCase.CreateAd(c.Request.Context(), &req); err != nil {
		logger.Error("建立廣告失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "建立廣告失敗",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "建立廣告成功",
	})
}

// GetAllAds 取得所有廣告（管理功能）
func (h *AdvertisementHandler) GetAllAds(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50 // 預設值
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0 // 預設值
	}

	ads, err := h.adUseCase.GetAllAds(c.Request.Context(), limit, offset)
	if err != nil {
		logger.Error("取得所有廣告失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "取得廣告失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"advertisements": ads,
		"count":          len(ads),
	})
}

// UpdateAd 更新廣告（管理功能）
func (h *AdvertisementHandler) UpdateAd(c *gin.Context) {
	idStr := c.Param("id")
	adID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的廣告 ID",
		})
		return
	}

	var req domain.CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新廣告請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	if err := h.adUseCase.UpdateAd(c.Request.Context(), adID, &req); err != nil {
		logger.Error("更新廣告失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新廣告失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新廣告成功",
	})
}

// DeleteAd 刪除廣告（管理功能）
func (h *AdvertisementHandler) DeleteAd(c *gin.Context) {
	idStr := c.Param("id")
	adID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的廣告 ID",
		})
		return
	}

	if err := h.adUseCase.DeleteAd(c.Request.Context(), adID); err != nil {
		logger.Error("刪除廣告失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "刪除廣告失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "刪除廣告成功",
	})
}

// RecordView 記錄廣告瀏覽
func (h *AdvertisementHandler) RecordView(c *gin.Context) {
	var req struct {
		AdvertisementID int    `json:"advertisement_id" binding:"required"`
		UserID          *int   `json:"user_id,omitempty"`
		GameSessionID   string `json:"game_session_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("記錄廣告瀏覽請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	// 處理可選的 UserID
	var userID int
	if req.UserID != nil {
		userID = *req.UserID
	}

	view := &domain.AdView{
		AdvertisementID: req.AdvertisementID,
		UserID:          userID,
		GameSessionID:   req.GameSessionID,
	}

	if err := h.adUseCase.RecordAdView(c.Request.Context(), view); err != nil {
		logger.Error("記錄廣告瀏覽失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "記錄瀏覽失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "記錄瀏覽成功",
	})
}

// RecordClick 記錄廣告點擊
func (h *AdvertisementHandler) RecordClick(c *gin.Context) {
	var req struct {
		AdvertisementID int    `json:"advertisement_id" binding:"required"`
		UserID          *int   `json:"user_id,omitempty"`
		GameSessionID   string `json:"game_session_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("記錄廣告點擊請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	// 處理可選的 UserID
	var userID int
	if req.UserID != nil {
		userID = *req.UserID
	}

	click := &domain.AdClick{
		AdvertisementID: req.AdvertisementID,
		UserID:          userID,
		GameSessionID:   req.GameSessionID,
	}

	if err := h.adUseCase.RecordAdClick(c.Request.Context(), click); err != nil {
		logger.Error("記錄廣告點擊失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "記錄點擊失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "記錄點擊成功",
	})
}

// GetAdByID 取得特定廣告
func (h *AdvertisementHandler) GetAdByID(c *gin.Context) {
	idStr := c.Param("id")
	adID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的廣告 ID",
		})
		return
	}

	ad, err := h.adUseCase.GetAdvertisementByID(c.Request.Context(), adID)
	if err != nil {
		logger.Error("取得廣告失敗", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "取得廣告失敗",
		})
		return
	}

	if ad == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "廣告不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"advertisement": ad,
	})
}
