package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// GameUseCase 遊戲業務邏輯
type GameUseCase struct {
	gameRepo       GameRepository
	restaurantRepo RestaurantRepository
	favoriteRepo   FavoriteRepository
	adRepo         AdvertisementRepository
}

// NewGameUseCase 建立遊戲用例
func NewGameUseCase(
	gameRepo GameRepository,
	restaurantRepo RestaurantRepository,
	favoriteRepo FavoriteRepository,
	adRepo AdvertisementRepository,
) *GameUseCase {
	return &GameUseCase{
		gameRepo:       gameRepo,
		restaurantRepo: restaurantRepo,
		favoriteRepo:   favoriteRepo,
		adRepo:         adRepo,
	}
}

// StartGame 開始遊戲
func (uc *GameUseCase) StartGame(ctx context.Context, userID int, req *domain.StartGameRequest) (*domain.GameSession, error) {
	// 產生遊戲會話 ID
	sessionID := uuid.New().String()

	// 取得附近餐廳
	searchParams := &domain.RestaurantSearchParams{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Radius:    req.Radius,
		Limit:     20, // 最多 20 間餐廳參與遊戲
	}

	nearbyRestaurants, err := uc.restaurantRepo.SearchNearby(ctx, searchParams)
	if err != nil {
		logger.Error("搜尋附近餐廳失敗", zap.Error(err))
		return nil, errors.New("搜尋餐廳失敗")
	}

	// 取得使用者最愛餐廳
	favorites, err := uc.favoriteRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Warn("取得最愛餐廳失敗", zap.Error(err), zap.Int("user_id", userID))
	}

	// 合併附近餐廳和最愛餐廳
	restaurants := uc.mergeRestaurants(nearbyRestaurants, favorites)

	// 確保至少有餐廳可以參與遊戲
	if len(restaurants) == 0 {
		return nil, errors.New("附近沒有找到餐廳")
	}

	// 取得活躍廣告
	advertisements, err := uc.adRepo.GetActiveAds(ctx, 3) // 最多 3 個廣告
	if err != nil {
		logger.Warn("取得廣告失敗", zap.Error(err))
		advertisements = []domain.Advertisement{} // 如果取得廣告失敗，繼續遊戲但沒有廣告
	}

	// 建立遊戲會話
	session := &domain.GameSession{
		ID:             sessionID,
		UserID:         userID,
		GameType:       req.GameType,
		Status:         "playing",
		Restaurants:    restaurants,
		Advertisements: advertisements,
		StartedAt:      time.Now(),
		CreatedAt:      time.Now(),
	}

	if err := uc.gameRepo.CreateSession(ctx, session); err != nil {
		logger.Error("建立遊戲會話失敗", zap.Error(err))
		return nil, errors.New("開始遊戲失敗")
	}

	// 記錄廣告瀏覽
	uc.recordAdViews(ctx, userID, sessionID, advertisements)

	logger.Info("遊戲開始", 
		zap.String("session_id", sessionID),
		zap.Int("user_id", userID),
		zap.String("game_type", string(req.GameType)),
		zap.Int("restaurant_count", len(restaurants)),
		zap.Int("ad_count", len(advertisements)),
	)

	return session, nil
}

// CompleteGame 完成遊戲
func (uc *GameUseCase) CompleteGame(ctx context.Context, userID int, req *domain.CompleteGameRequest) (*domain.GameResult, error) {
	// 取得遊戲會話
	session, err := uc.gameRepo.GetSessionByID(ctx, req.SessionID)
	if err != nil {
		logger.Error("取得遊戲會話失敗", zap.Error(err))
		return nil, errors.New("遊戲會話不存在")
	}

	// 驗證使用者權限
	if session.UserID != userID {
		return nil, errors.New("沒有權限操作此遊戲")
	}

	// 檢查遊戲狀態
	if session.Status != "playing" {
		return nil, errors.New("遊戲已結束")
	}

	// 找到選中的餐廳
	var selectedRestaurant *domain.RestaurantWithDistance
	for _, restaurant := range session.Restaurants {
		if restaurant.ID == req.SelectedRestaurantID {
			selectedRestaurant = &restaurant
			break
		}
	}

	if selectedRestaurant == nil {
		return nil, errors.New("選中的餐廳不在遊戲列表中")
	}

	// 處理廣告點擊
	var clickedAd *domain.Advertisement
	if req.ClickedAdID != nil {
		for _, ad := range session.Advertisements {
			if ad.ID == *req.ClickedAdID {
				clickedAd = &ad
				// 記錄廣告點擊
				uc.recordAdClick(ctx, userID, session.ID, ad.ID)
				// 如果點擊了廣告，用廣告餐廳替換選中的餐廳
				if adRestaurant, err := uc.restaurantRepo.GetByID(ctx, ad.RestaurantID); err == nil {
					selectedRestaurant = &domain.RestaurantWithDistance{
						Restaurant: *adRestaurant,
						Distance:   0, // 廣告餐廳距離設為 0
					}
				}
				break
			}
		}
	}

	// 更新遊戲會話
	completedAt := time.Now()
	session.Status = "completed"
	session.Result = selectedRestaurant
	session.CompletedAt = &completedAt

	if err := uc.gameRepo.UpdateSession(ctx, session); err != nil {
		logger.Error("更新遊戲會話失敗", zap.Error(err))
		return nil, errors.New("完成遊戲失敗")
	}

	result := &domain.GameResult{
		SessionID:          req.SessionID,
		SelectedRestaurant: selectedRestaurant,
		ClickedAd:          clickedAd,
		CompletedAt:        completedAt,
	}

	logger.Info("遊戲完成",
		zap.String("session_id", req.SessionID),
		zap.Int("user_id", userID),
		zap.Int("selected_restaurant_id", selectedRestaurant.ID),
		zap.Bool("clicked_ad", clickedAd != nil),
	)

	return result, nil
}

// GetGameHistory 取得遊戲歷史
func (uc *GameUseCase) GetGameHistory(ctx context.Context, userID int, limit, offset int) ([]domain.GameSession, error) {
	sessions, err := uc.gameRepo.GetUserSessions(ctx, userID, limit, offset)
	if err != nil {
		logger.Error("取得遊戲歷史失敗", zap.Error(err), zap.Int("user_id", userID))
		return nil, errors.New("取得遊戲歷史失敗")
	}

	return sessions, nil
}

// mergeRestaurants 合併附近餐廳和最愛餐廳
func (uc *GameUseCase) mergeRestaurants(nearby []domain.RestaurantWithDistance, favorites []domain.FavoriteRestaurant) []domain.RestaurantWithDistance {
	restaurantMap := make(map[int]domain.RestaurantWithDistance)

	// 先加入附近餐廳
	for _, restaurant := range nearby {
		restaurantMap[restaurant.ID] = restaurant
	}

	// 加入最愛餐廳（如果不在附近餐廳中）
	for _, fav := range favorites {
		if _, exists := restaurantMap[fav.RestaurantID]; !exists {
			// 這裡需要取得餐廳詳細資訊，暫時跳過
			// 在實際實作中，可以呼叫 restaurantRepo.GetByID
		}
	}

	// 轉換為切片並隨機排序
	restaurants := make([]domain.RestaurantWithDistance, 0, len(restaurantMap))
	for _, restaurant := range restaurantMap {
		restaurants = append(restaurants, restaurant)
	}

	// 隨機打亂順序
	rand.Shuffle(len(restaurants), func(i, j int) {
		restaurants[i], restaurants[j] = restaurants[j], restaurants[i]
	})

	return restaurants
}

// recordAdViews 記錄廣告瀏覽
func (uc *GameUseCase) recordAdViews(ctx context.Context, userID int, sessionID string, ads []domain.Advertisement) {
	for _, ad := range ads {
		view := &domain.AdView{
			AdvertisementID: ad.ID,
			UserID:          userID,
			GameSessionID:   sessionID,
			ViewedAt:        time.Now(),
		}
		if err := uc.adRepo.RecordView(ctx, view); err != nil {
			logger.Warn("記錄廣告瀏覽失敗", zap.Error(err), zap.Int("ad_id", ad.ID))
		}
	}
}

// recordAdClick 記錄廣告點擊
func (uc *GameUseCase) recordAdClick(ctx context.Context, userID int, sessionID string, adID int) {
	click := &domain.AdClick{
		AdvertisementID: adID,
		UserID:          userID,
		GameSessionID:   sessionID,
		ClickedAt:       time.Now(),
	}
	if err := uc.adRepo.RecordClick(ctx, click); err != nil {
		logger.Warn("記錄廣告點擊失敗", zap.Error(err), zap.Int("ad_id", adID))
	}
}