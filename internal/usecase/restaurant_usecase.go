package usecase

import (
	"context"
	"errors"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// RestaurantUseCase 餐廳業務邏輯
type RestaurantUseCase struct {
	restaurantRepo RestaurantRepository
	favoriteRepo   FavoriteRepository
	externalAPI    ExternalAPIService
}

// NewRestaurantUseCase 建立餐廳用例
func NewRestaurantUseCase(
	restaurantRepo RestaurantRepository,
	favoriteRepo FavoriteRepository,
	externalAPI ExternalAPIService,
) *RestaurantUseCase {
	return &RestaurantUseCase{
		restaurantRepo: restaurantRepo,
		favoriteRepo:   favoriteRepo,
		externalAPI:    externalAPI,
	}
}

// SearchNearby 搜尋附近餐廳
func (uc *RestaurantUseCase) SearchNearby(ctx context.Context, params *domain.RestaurantSearchParams) ([]domain.RestaurantWithDistance, error) {
	// 先從本地資料庫搜尋
	restaurants, err := uc.restaurantRepo.SearchNearby(ctx, params)
	if err != nil {
		logger.Error("搜尋本地餐廳失敗", zap.Error(err))
		return nil, errors.New("搜尋餐廳失敗")
	}

	// 如果本地餐廳數量不足，可以從外部 API 補充
	if len(restaurants) < 5 && uc.externalAPI != nil {
		externalRestaurants, err := uc.externalAPI.SearchNearbyRestaurants(ctx, params.Latitude, params.Longitude, params.Radius)
		if err != nil {
			logger.Warn("從外部 API 搜尋餐廳失敗", zap.Error(err))
		} else {
			// 將外部餐廳加入本地資料庫
			for _, restaurant := range externalRestaurants {
				if err := uc.restaurantRepo.Create(ctx, &restaurant); err != nil {
					logger.Warn("儲存外部餐廳失敗", zap.Error(err), zap.String("name", restaurant.Name))
				}
			}
			
			// 重新搜尋
			restaurants, err = uc.restaurantRepo.SearchNearby(ctx, params)
			if err != nil {
				logger.Error("重新搜尋餐廳失敗", zap.Error(err))
			}
		}
	}

	logger.Info("搜尋附近餐廳完成", 
		zap.Float64("latitude", params.Latitude),
		zap.Float64("longitude", params.Longitude),
		zap.Int("radius", params.Radius),
		zap.Int("count", len(restaurants)),
	)

	return restaurants, nil
}

// GetRestaurant 取得餐廳詳細資訊
func (uc *RestaurantUseCase) GetRestaurant(ctx context.Context, id int) (*domain.Restaurant, error) {
	restaurant, err := uc.restaurantRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("取得餐廳資訊失敗", zap.Error(err), zap.Int("restaurant_id", id))
		return nil, errors.New("餐廳不存在")
	}

	return restaurant, nil
}

// AddToFavorites 新增到最愛餐廳
func (uc *RestaurantUseCase) AddToFavorites(ctx context.Context, userID int, req *domain.AddFavoriteRequest) error {
	// 檢查餐廳是否存在
	_, err := uc.restaurantRepo.GetByID(ctx, req.RestaurantID)
	if err != nil {
		logger.Error("餐廳不存在", zap.Error(err), zap.Int("restaurant_id", req.RestaurantID))
		return errors.New("餐廳不存在")
	}

	// 檢查是否已經在最愛中
	exists, err := uc.favoriteRepo.IsExists(ctx, userID, req.RestaurantID)
	if err != nil {
		logger.Error("檢查最愛餐廳失敗", zap.Error(err))
		return errors.New("操作失敗")
	}

	if exists {
		return errors.New("餐廳已在最愛清單中")
	}

	// 新增到最愛
	if err := uc.favoriteRepo.Add(ctx, userID, req); err != nil {
		logger.Error("新增最愛餐廳失敗", zap.Error(err), zap.Int("user_id", userID), zap.Int("restaurant_id", req.RestaurantID))
		return errors.New("新增最愛失敗")
	}

	logger.Info("新增最愛餐廳成功", zap.Int("user_id", userID), zap.Int("restaurant_id", req.RestaurantID))
	return nil
}

// RemoveFromFavorites 從最愛餐廳移除
func (uc *RestaurantUseCase) RemoveFromFavorites(ctx context.Context, userID, restaurantID int) error {
	// 檢查是否在最愛中
	exists, err := uc.favoriteRepo.IsExists(ctx, userID, restaurantID)
	if err != nil {
		logger.Error("檢查最愛餐廳失敗", zap.Error(err))
		return errors.New("操作失敗")
	}

	if !exists {
		return errors.New("餐廳不在最愛清單中")
	}

	// 從最愛移除
	if err := uc.favoriteRepo.Remove(ctx, userID, restaurantID); err != nil {
		logger.Error("移除最愛餐廳失敗", zap.Error(err), zap.Int("user_id", userID), zap.Int("restaurant_id", restaurantID))
		return errors.New("移除最愛失敗")
	}

	logger.Info("移除最愛餐廳成功", zap.Int("user_id", userID), zap.Int("restaurant_id", restaurantID))
	return nil
}

// GetFavorites 取得使用者最愛餐廳
func (uc *RestaurantUseCase) GetFavorites(ctx context.Context, userID int) ([]domain.FavoriteRestaurant, error) {
	favorites, err := uc.favoriteRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("取得最愛餐廳失敗", zap.Error(err), zap.Int("user_id", userID))
		return nil, errors.New("取得最愛餐廳失敗")
	}

	logger.Info("取得最愛餐廳成功", zap.Int("user_id", userID), zap.Int("count", len(favorites)))
	return favorites, nil
}

// CreateRestaurant 建立餐廳（管理功能）
func (uc *RestaurantUseCase) CreateRestaurant(ctx context.Context, restaurant *domain.Restaurant) error {
	if err := uc.restaurantRepo.Create(ctx, restaurant); err != nil {
		logger.Error("建立餐廳失敗", zap.Error(err), zap.String("name", restaurant.Name))
		return errors.New("建立餐廳失敗")
	}

	logger.Info("建立餐廳成功", zap.String("name", restaurant.Name), zap.Int("id", restaurant.ID))
	return nil
}

// GetAllRestaurants 取得所有餐廳（管理功能）
func (uc *RestaurantUseCase) GetAllRestaurants(ctx context.Context, limit, offset int) ([]domain.Restaurant, error) {
	restaurants, err := uc.restaurantRepo.GetAll(ctx, limit, offset)
	if err != nil {
		logger.Error("取得餐廳清單失敗", zap.Error(err))
		return nil, errors.New("取得餐廳清單失敗")
	}

	return restaurants, nil
}