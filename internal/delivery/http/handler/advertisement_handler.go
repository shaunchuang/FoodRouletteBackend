package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// AdvertisementHandler 廣告 HTTP 處理器
type AdvertisementHandler struct {
	// TODO: 加入廣告用例
	// adUseCase *usecase.AdvertisementUseCase
}

// NewAdvertisementHandler 建立廣告處理器
func NewAdvertisementHandler() *AdvertisementHandler {
	return &AdvertisementHandler{
		// adUseCase: adUseCase,
	}
}

// GetActiveAds 取得活躍廣告
func (h *AdvertisementHandler) GetActiveAds(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	_, err := strconv.Atoi(limitStr)
	if err != nil {
		// 使用預設值
	}

	// TODO: 實作取得活躍廣告邏輯
	// ads, err := h.adUseCase.GetActiveAds(c.Request.Context(), limit)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": err.Error(),
	//     })
	//     return
	// }

	// 暫時返回空陣列
	c.JSON(http.StatusOK, gin.H{
		"advertisements": []domain.Advertisement{},
		"count":          0,
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

	// TODO: 實作取得廣告統計邏輯
	// stats, err := h.adUseCase.GetStatistics(c.Request.Context(), adID, period)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": err.Error(),
	//     })
	//     return
	// }

	// 暫時返回空統計
	stats := &domain.AdStatistics{
		AdvertisementID: adID,
		ViewCount:       0,
		ClickCount:      0,
		CTR:             0.0,
		Period:          period,
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
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	// TODO: 實作建立廣告邏輯
	// if err := h.adUseCase.CreateAd(c.Request.Context(), &req); err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": err.Error(),
	//     })
	//     return
	// }

	c.JSON(http.StatusCreated, gin.H{
		"message": "建立廣告成功",
	})
}

// GetAllAds 取得所有廣告（管理功能）
func (h *AdvertisementHandler) GetAllAds(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	_, err := strconv.Atoi(limitStr)
	if err != nil {
		// 使用預設值
	}

	_, err = strconv.Atoi(offsetStr)
	if err != nil {
		// 使用預設值
	}

	// TODO: 實作取得所有廣告邏輯
	// ads, err := h.adUseCase.GetAllAds(c.Request.Context(), limit, offset)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": err.Error(),
	//     })
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"advertisements": []domain.Advertisement{},
		"count":          0,
	})
}

// UpdateAd 更新廣告（管理功能）
func (h *AdvertisementHandler) UpdateAd(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.Atoi(idStr)
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
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	// TODO: 實作更新廣告邏輯
	// if err := h.adUseCase.UpdateAd(c.Request.Context(), adID, &req); err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": err.Error(),
	//     })
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "更新廣告成功",
	})
}

// DeleteAd 刪除廣告（管理功能）
func (h *AdvertisementHandler) DeleteAd(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的廣告 ID",
		})
		return
	}

	// TODO: 實作刪除廣告邏輯
	// if err := h.adUseCase.DeleteAd(c.Request.Context(), adID); err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": err.Error(),
	//     })
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "刪除廣告成功",
	})
}