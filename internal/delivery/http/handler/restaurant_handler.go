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

// RestaurantHandler 餐廳 HTTP 處理器
type RestaurantHandler struct {
	restaurantUseCase *usecase.RestaurantUseCase
}

// NewRestaurantHandler 建立餐廳處理器
func NewRestaurantHandler(restaurantUseCase *usecase.RestaurantUseCase) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantUseCase: restaurantUseCase,
	}
}

// SearchNearby 搜尋附近餐廳
func (h *RestaurantHandler) SearchNearby(c *gin.Context) {
	var params domain.RestaurantSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Error("搜尋餐廳請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	// 設定預設值
	if params.Radius == 0 {
		params.Radius = 1000 // 預設 1 公里
	}
	if params.Limit == 0 {
		params.Limit = 20 // 預設最多 20 間
	}

	restaurants, err := h.restaurantUseCase.SearchNearby(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurants": restaurants,
		"count":       len(restaurants),
	})
}

// GetRestaurant 取得餐廳詳細資訊
func (h *RestaurantHandler) GetRestaurant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的餐廳 ID",
		})
		return
	}

	restaurant, err := h.restaurantUseCase.GetRestaurant(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant": restaurant,
	})
}

// AddToFavorites 新增到最愛餐廳
func (h *RestaurantHandler) AddToFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	var req domain.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("新增最愛餐廳請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	if err := h.restaurantUseCase.AddToFavorites(c.Request.Context(), userID.(int), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "新增最愛餐廳成功",
	})
}

// RemoveFromFavorites 從最愛餐廳移除
func (h *RestaurantHandler) RemoveFromFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	restaurantIDStr := c.Param("restaurant_id")
	restaurantID, err := strconv.Atoi(restaurantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的餐廳 ID",
		})
		return
	}

	if err := h.restaurantUseCase.RemoveFromFavorites(c.Request.Context(), userID.(int), restaurantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "移除最愛餐廳成功",
	})
}

// GetFavorites 取得使用者最愛餐廳
func (h *RestaurantHandler) GetFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的使用者",
		})
		return
	}

	favorites, err := h.restaurantUseCase.GetFavorites(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"favorites": favorites,
		"count":     len(favorites),
	})
}

// CreateRestaurant 建立餐廳（管理功能）
func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	var restaurant domain.Restaurant
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		logger.Error("建立餐廳請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	if err := h.restaurantUseCase.CreateRestaurant(c.Request.Context(), &restaurant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "建立餐廳成功",
		"restaurant": restaurant,
	})
}

// GetAllRestaurants 取得所有餐廳（管理功能）
func (h *RestaurantHandler) GetAllRestaurants(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	restaurants, err := h.restaurantUseCase.GetAllRestaurants(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurants": restaurants,
		"count":       len(restaurants),
	})
}

// UpdateRestaurant 更新餐廳資訊（管理功能）
func (h *RestaurantHandler) UpdateRestaurant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的餐廳 ID",
		})
		return
	}

	var restaurant domain.Restaurant
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		logger.Error("更新餐廳請求參數錯誤", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤",
			"details": err.Error(),
		})
		return
	}

	restaurant.ID = id
	// TODO: 實作更新餐廳邏輯
	// if err := h.restaurantUseCase.UpdateRestaurant(c.Request.Context(), &restaurant); err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": err.Error(),
	//     })
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "更新餐廳成功",
		"restaurant": restaurant,
	})
}