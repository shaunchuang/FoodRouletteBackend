package usecase

import (
	"context"
	"time"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// AdvertisementUseCase 廣告業務邏輯
type AdvertisementUseCase struct {
	adRepo AdvertisementRepository
}

// NewAdvertisementUseCase 建立廣告 Use Case
func NewAdvertisementUseCase(adRepo AdvertisementRepository) *AdvertisementUseCase {
	return &AdvertisementUseCase{
		adRepo: adRepo,
	}
}

// GetActiveAdvertisements 取得活躍的廣告
func (uc *AdvertisementUseCase) GetActiveAdvertisements(ctx context.Context, limit int) ([]domain.Advertisement, error) {
	logger.Info("取得活躍廣告", zap.Int("limit", limit))

	ads, err := uc.adRepo.GetActiveAds(ctx, limit)
	if err != nil {
		logger.Error("取得活躍廣告失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("取得活躍廣告成功", zap.Int("count", len(ads)))
	return ads, nil
}

// GetAdvertisementByID 根據 ID 取得廣告
func (uc *AdvertisementUseCase) GetAdvertisementByID(ctx context.Context, id int) (*domain.Advertisement, error) {
	logger.Info("取得廣告詳細資訊", zap.Int("ad_id", id))

	ad, err := uc.adRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("取得廣告失敗", zap.Error(err), zap.Int("ad_id", id))
		return nil, err
	}

	logger.Info("取得廣告成功", zap.Int("ad_id", id))
	return ad, nil
}

// CreateAdvertisement 建立新廣告
func (uc *AdvertisementUseCase) CreateAdvertisement(ctx context.Context, ad *domain.Advertisement) error {
	logger.Info("建立新廣告", zap.String("title", ad.Title))

	// 設定建立時間
	now := time.Now()
	ad.CreatedAt = now
	ad.UpdatedAt = now

	// 初始化計數器
	ad.ClickCount = 0
	ad.ViewCount = 0

	if err := uc.adRepo.Create(ctx, ad); err != nil {
		logger.Error("建立廣告失敗", zap.Error(err), zap.String("title", ad.Title))
		return err
	}

	logger.Info("廣告建立成功", zap.Int("ad_id", ad.ID), zap.String("title", ad.Title))
	return nil
}

// UpdateAdvertisement 更新廣告
func (uc *AdvertisementUseCase) UpdateAdvertisement(ctx context.Context, ad *domain.Advertisement) error {
	logger.Info("更新廣告", zap.Int("ad_id", ad.ID))

	// 檢查廣告是否存在
	_, err := uc.adRepo.GetByID(ctx, ad.ID)
	if err != nil {
		logger.Error("廣告不存在", zap.Error(err), zap.Int("ad_id", ad.ID))
		return err
	}

	if err := uc.adRepo.Update(ctx, ad); err != nil {
		logger.Error("更新廣告失敗", zap.Error(err), zap.Int("ad_id", ad.ID))
		return err
	}

	logger.Info("廣告更新成功", zap.Int("ad_id", ad.ID))
	return nil
}

// RecordAdView 記錄廣告瀏覽
func (uc *AdvertisementUseCase) RecordAdView(ctx context.Context, view *domain.AdView) error {
	logger.Info("記錄廣告瀏覽",
		zap.Int("ad_id", view.AdvertisementID),
		zap.Int("user_id", view.UserID),
		zap.String("session_id", view.GameSessionID),
	)

	// 設定瀏覽時間
	view.ViewedAt = time.Now()

	if err := uc.adRepo.RecordView(ctx, view); err != nil {
		logger.Error("記錄廣告瀏覽失敗", zap.Error(err))
		return err
	}

	logger.Info("廣告瀏覽記錄成功")
	return nil
}

// RecordAdClick 記錄廣告點擊
func (uc *AdvertisementUseCase) RecordAdClick(ctx context.Context, click *domain.AdClick) error {
	logger.Info("記錄廣告點擊",
		zap.Int("ad_id", click.AdvertisementID),
		zap.Int("user_id", click.UserID),
		zap.String("session_id", click.GameSessionID),
	)

	// 設定點擊時間
	click.ClickedAt = time.Now()

	if err := uc.adRepo.RecordClick(ctx, click); err != nil {
		logger.Error("記錄廣告點擊失敗", zap.Error(err))
		return err
	}

	logger.Info("廣告點擊記錄成功")
	return nil
}

// GetAdvertisementStatistics 取得廣告統計資訊
func (uc *AdvertisementUseCase) GetAdvertisementStatistics(ctx context.Context, adID int, period string) (*domain.AdStatistics, error) {
	logger.Info("取得廣告統計", zap.Int("ad_id", adID), zap.String("period", period))

	// 驗證時間週期參數
	validPeriods := map[string]bool{
		"day":   true,
		"week":  true,
		"month": true,
		"all":   true,
	}

	if !validPeriods[period] {
		logger.Error("無效的時間週期", zap.String("period", period))
		return nil, domain.ErrInvalidPeriod
	}

	stats, err := uc.adRepo.GetStatistics(ctx, adID, period)
	if err != nil {
		logger.Error("取得廣告統計失敗", zap.Error(err), zap.Int("ad_id", adID))
		return nil, err
	}

	logger.Info("廣告統計取得成功",
		zap.Int("ad_id", adID),
		zap.Int64("views", int64(stats.ViewCount)),
		zap.Int64("clicks", int64(stats.ClickCount)),
	)

	return stats, nil
}

// GetAdvertisementsForGame 為遊戲取得適合的廣告
func (uc *AdvertisementUseCase) GetAdvertisementsForGame(ctx context.Context, userID int, gameType string) ([]domain.Advertisement, error) {
	logger.Info("為遊戲取得廣告",
		zap.Int("user_id", userID),
		zap.String("game_type", gameType),
	)

	// 取得活躍廣告，可以根據遊戲類型和使用者偏好進行過濾
	// 這裡先簡單實作，未來可以加入更複雜的推薦演算法
	ads, err := uc.adRepo.GetActiveAds(ctx, 5) // 限制返回 5 個廣告
	if err != nil {
		logger.Error("取得遊戲廣告失敗", zap.Error(err))
		return nil, err
	}

	// 可以在這裡加入廣告排序邏輯，例如根據優先級、點擊率等

	logger.Info("遊戲廣告取得成功",
		zap.Int("count", len(ads)),
		zap.Int("user_id", userID),
	)

	return ads, nil
}

// CreateAd 建立新廣告方法
func (uc *AdvertisementUseCase) CreateAd(ctx context.Context, req *domain.CreateAdRequest) error {
	logger.Info("建立新廣告", zap.String("title", req.Title))

	// 建立廣告物件
	ad := &domain.Advertisement{
		RestaurantID: req.RestaurantID,
		Title:        req.Title,
		Content:      req.Content,
		ImageURL:     req.ImageURL,
		TargetURL:    req.TargetURL,
		IsActive:     true, // 預設為活躍
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Priority:     req.Priority,
		ClickCount:   0,
		ViewCount:    0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.adRepo.Create(ctx, ad); err != nil {
		logger.Error("建立廣告失敗", zap.Error(err))
		return err
	}

	logger.Info("廣告建立成功", zap.String("title", req.Title))
	return nil
}

// UpdateAd 更新廣告方法
func (uc *AdvertisementUseCase) UpdateAd(ctx context.Context, adID int, req *domain.CreateAdRequest) error {
	logger.Info("更新廣告", zap.Int("ad_id", adID))

	// 檢查廣告是否存在
	existingAd, err := uc.adRepo.GetByID(ctx, adID)
	if err != nil {
		logger.Error("廣告不存在", zap.Error(err), zap.Int("ad_id", adID))
		return err
	}

	// 更新廣告資料
	existingAd.RestaurantID = req.RestaurantID
	existingAd.Title = req.Title
	existingAd.Content = req.Content
	existingAd.ImageURL = req.ImageURL
	existingAd.TargetURL = req.TargetURL
	existingAd.StartDate = req.StartDate
	existingAd.EndDate = req.EndDate
	existingAd.Priority = req.Priority
	existingAd.UpdatedAt = time.Now()

	if err := uc.adRepo.Update(ctx, existingAd); err != nil {
		logger.Error("更新廣告失敗", zap.Error(err), zap.Int("ad_id", adID))
		return err
	}

	logger.Info("廣告更新成功", zap.Int("ad_id", adID))
	return nil
}

// DeleteAd 刪除廣告方法
func (uc *AdvertisementUseCase) DeleteAd(ctx context.Context, adID int) error {
	logger.Info("刪除廣告", zap.Int("ad_id", adID))

	// 檢查廣告是否存在
	_, err := uc.adRepo.GetByID(ctx, adID)
	if err != nil {
		logger.Error("廣告不存在", zap.Error(err), zap.Int("ad_id", adID))
		return err
	}

	// 軟刪除：將 IsActive 設為 false
	ad := &domain.Advertisement{
		ID:        adID,
		IsActive:  false,
		UpdatedAt: time.Now(),
	}

	if err := uc.adRepo.Update(ctx, ad); err != nil {
		logger.Error("刪除廣告失敗", zap.Error(err), zap.Int("ad_id", adID))
		return err
	}

	logger.Info("廣告刪除成功", zap.Int("ad_id", adID))
	return nil
}

// GetAllAds 取得所有廣告方法
func (uc *AdvertisementUseCase) GetAllAds(ctx context.Context, limit, offset int) ([]domain.Advertisement, error) {
	logger.Info("取得所有廣告", zap.Int("limit", limit), zap.Int("offset", offset))

	// 注意：這裡需要在 repository interface 新增 GetAll 方法
	// 現在先用 GetActiveAds，但應該要實作 GetAll 方法包含非活躍廣告
	ads, err := uc.adRepo.GetActiveAds(ctx, limit)
	if err != nil {
		logger.Error("取得所有廣告失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("取得所有廣告成功", zap.Int("count", len(ads)))
	return ads, nil
}
