package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// AdvertisementRepository PostgreSQL 廣告資料庫操作實作
type AdvertisementRepository struct {
	db *sql.DB
}

// NewAdvertisementRepository 建立廣告 Repository
func NewAdvertisementRepository(db *sql.DB) *AdvertisementRepository {
	return &AdvertisementRepository{
		db: db,
	}
}

// Create 建立新廣告
func (r *AdvertisementRepository) Create(ctx context.Context, ad *domain.Advertisement) error {
	query := `
		INSERT INTO advertisements (restaurant_id, title, content, image_url, target_url, is_active, start_date, end_date, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		ad.RestaurantID,
		ad.Title,
		ad.Content,
		ad.ImageURL,
		ad.TargetURL,
		ad.IsActive,
		ad.StartDate,
		ad.EndDate,
		ad.Priority,
		now,
		now,
	).Scan(&ad.ID)

	if err != nil {
		logger.Error("建立廣告失敗", zap.Error(err), zap.String("title", ad.Title))
		return err
	}

	ad.CreatedAt = now
	ad.UpdatedAt = now

	logger.Info("廣告建立成功", zap.Int("ad_id", ad.ID), zap.String("title", ad.Title))
	return nil
}

// GetByID 根據 ID 取得廣告
func (r *AdvertisementRepository) GetByID(ctx context.Context, id int) (*domain.Advertisement, error) {
	query := `
		SELECT id, restaurant_id, title, content, image_url, target_url, is_active, 
		       start_date, end_date, priority, click_count, view_count, created_at, updated_at
		FROM advertisements
		WHERE id = $1`

	ad := &domain.Advertisement{}
	var imageURL, targetURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ad.ID,
		&ad.RestaurantID,
		&ad.Title,
		&ad.Content,
		&imageURL,
		&targetURL,
		&ad.IsActive,
		&ad.StartDate,
		&ad.EndDate,
		&ad.Priority,
		&ad.ClickCount,
		&ad.ViewCount,
		&ad.CreatedAt,
		&ad.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("廣告不存在")
		}
		logger.Error("取得廣告失敗", zap.Error(err), zap.Int("ad_id", id))
		return nil, err
	}

	// 處理可為空的欄位
	if imageURL.Valid {
		ad.ImageURL = imageURL.String
	}
	if targetURL.Valid {
		ad.TargetURL = targetURL.String
	}

	return ad, nil
}

// GetActiveAds 取得活躍的廣告
func (r *AdvertisementRepository) GetActiveAds(ctx context.Context, limit int) ([]domain.Advertisement, error) {
	query := `
		SELECT id, restaurant_id, title, content, image_url, target_url, is_active, 
		       start_date, end_date, priority, click_count, view_count, created_at, updated_at
		FROM advertisements
		WHERE is_active = TRUE 
		  AND start_date <= CURRENT_TIMESTAMP 
		  AND end_date >= CURRENT_TIMESTAMP
		ORDER BY priority DESC, created_at DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		logger.Error("取得活躍廣告失敗", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var ads []domain.Advertisement
	for rows.Next() {
		var ad domain.Advertisement
		var imageURL, targetURL sql.NullString

		err := rows.Scan(
			&ad.ID,
			&ad.RestaurantID,
			&ad.Title,
			&ad.Content,
			&imageURL,
			&targetURL,
			&ad.IsActive,
			&ad.StartDate,
			&ad.EndDate,
			&ad.Priority,
			&ad.ClickCount,
			&ad.ViewCount,
			&ad.CreatedAt,
			&ad.UpdatedAt,
		)

		if err != nil {
			logger.Error("掃描廣告資料失敗", zap.Error(err))
			continue
		}

		// 處理可為空的欄位
		if imageURL.Valid {
			ad.ImageURL = imageURL.String
		}
		if targetURL.Valid {
			ad.TargetURL = targetURL.String
		}

		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		logger.Error("處理廣告查詢結果失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("取得活躍廣告成功", zap.Int("count", len(ads)))
	return ads, nil
}

// GetAll 取得所有廣告（包含非活躍的）
func (r *AdvertisementRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Advertisement, error) {
	query := `
		SELECT id, restaurant_id, title, content, image_url, target_url, 
		       is_active, start_date, end_date, priority, click_count, view_count,
		       created_at, updated_at
		FROM advertisements
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		logger.Error("查詢所有廣告失敗", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var advertisements []domain.Advertisement
	for rows.Next() {
		var ad domain.Advertisement
		var imageURL, targetURL sql.NullString

		err := rows.Scan(
			&ad.ID,
			&ad.RestaurantID,
			&ad.Title,
			&ad.Content,
			&imageURL,
			&targetURL,
			&ad.IsActive,
			&ad.StartDate,
			&ad.EndDate,
			&ad.Priority,
			&ad.ClickCount,
			&ad.ViewCount,
			&ad.CreatedAt,
			&ad.UpdatedAt,
		)

		if err != nil {
			logger.Error("掃描廣告資料失敗", zap.Error(err))
			return nil, err
		}

		// 處理可為空的欄位
		if imageURL.Valid {
			ad.ImageURL = imageURL.String
		}
		if targetURL.Valid {
			ad.TargetURL = targetURL.String
		}

		advertisements = append(advertisements, ad)
	}

	if err = rows.Err(); err != nil {
		logger.Error("遍歷廣告資料失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("取得所有廣告成功", zap.Int("count", len(advertisements)))
	return advertisements, nil
}

// Update 更新廣告
func (r *AdvertisementRepository) Update(ctx context.Context, ad *domain.Advertisement) error {
	query := `
		UPDATE advertisements
		SET restaurant_id = $1, title = $2, content = $3, image_url = $4, target_url = $5,
		    is_active = $6, start_date = $7, end_date = $8, priority = $9, updated_at = $10
		WHERE id = $11`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		ad.RestaurantID,
		ad.Title,
		ad.Content,
		ad.ImageURL,
		ad.TargetURL,
		ad.IsActive,
		ad.StartDate,
		ad.EndDate,
		ad.Priority,
		now,
		ad.ID,
	)

	if err != nil {
		logger.Error("更新廣告失敗", zap.Error(err), zap.Int("ad_id", ad.ID))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("廣告不存在")
	}

	ad.UpdatedAt = now
	logger.Info("廣告更新成功", zap.Int("ad_id", ad.ID))
	return nil
}

// RecordView 記錄廣告瀏覽
func (r *AdvertisementRepository) RecordView(ctx context.Context, view *domain.AdView) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入瀏覽記錄
	insertQuery := `
		INSERT INTO ad_views (advertisement_id, user_id, game_session_id, viewed_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, insertQuery,
		view.AdvertisementID,
		view.UserID,
		view.GameSessionID,
		view.ViewedAt,
		view.IPAddress,
		view.UserAgent,
	)

	if err != nil {
		logger.Error("插入廣告瀏覽記錄失敗", zap.Error(err))
		return err
	}

	// 更新廣告瀏覽計數
	updateQuery := `
		UPDATE advertisements
		SET view_count = view_count + 1
		WHERE id = $1`

	_, err = tx.ExecContext(ctx, updateQuery, view.AdvertisementID)
	if err != nil {
		logger.Error("更新廣告瀏覽計數失敗", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.Error("提交廣告瀏覽記錄失敗", zap.Error(err))
		return err
	}

	logger.Info("記錄廣告瀏覽成功", zap.Int("ad_id", view.AdvertisementID), zap.Int("user_id", view.UserID))
	return nil
}

// RecordClick 記錄廣告點擊
func (r *AdvertisementRepository) RecordClick(ctx context.Context, click *domain.AdClick) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入點擊記錄
	insertQuery := `
		INSERT INTO ad_clicks (advertisement_id, user_id, game_session_id, clicked_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, insertQuery,
		click.AdvertisementID,
		click.UserID,
		click.GameSessionID,
		click.ClickedAt,
		click.IPAddress,
		click.UserAgent,
	)

	if err != nil {
		logger.Error("插入廣告點擊記錄失敗", zap.Error(err))
		return err
	}

	// 更新廣告點擊計數
	updateQuery := `
		UPDATE advertisements
		SET click_count = click_count + 1
		WHERE id = $1`

	_, err = tx.ExecContext(ctx, updateQuery, click.AdvertisementID)
	if err != nil {
		logger.Error("更新廣告點擊計數失敗", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.Error("提交廣告點擊記錄失敗", zap.Error(err))
		return err
	}

	logger.Info("記錄廣告點擊成功", zap.Int("ad_id", click.AdvertisementID), zap.Int("user_id", click.UserID))
	return nil
}

// GetStatistics 取得廣告統計資訊
func (r *AdvertisementRepository) GetStatistics(ctx context.Context, adID int, period string) (*domain.AdStatistics, error) {
	var dateFilter string
	switch period {
	case "day":
		dateFilter = "WHERE viewed_at >= CURRENT_DATE"
	case "week":
		dateFilter = "WHERE viewed_at >= CURRENT_DATE - INTERVAL '7 days'"
	case "month":
		dateFilter = "WHERE viewed_at >= CURRENT_DATE - INTERVAL '30 days'"
	default:
		dateFilter = "" // 全部時間
	}

	// 取得瀏覽統計
	viewQuery := `
		SELECT COUNT(*) as view_count,
		       COUNT(DISTINCT user_id) as unique_viewers
		FROM ad_views
		WHERE advertisement_id = $1 ` + dateFilter

	var viewCount, uniqueViewers int64
	err := r.db.QueryRowContext(ctx, viewQuery, adID).Scan(&viewCount, &uniqueViewers)
	if err != nil {
		logger.Error("取得廣告瀏覽統計失敗", zap.Error(err), zap.Int("ad_id", adID))
		return nil, err
	}

	// 取得點擊統計
	clickQuery := `
		SELECT COUNT(*) as click_count,
		       COUNT(DISTINCT user_id) as unique_clickers
		FROM ad_clicks
		WHERE advertisement_id = $1 ` + dateFilter

	var clickCount, uniqueClickers int64
	err = r.db.QueryRowContext(ctx, clickQuery, adID).Scan(&clickCount, &uniqueClickers)
	if err != nil {
		logger.Error("取得廣告點擊統計失敗", zap.Error(err), zap.Int("ad_id", adID))
		return nil, err
	}

	// 計算點擊率
	var ctr float64
	if viewCount > 0 {
		ctr = float64(clickCount) / float64(viewCount) * 100
	}

	stats := &domain.AdStatistics{
		AdvertisementID: adID,
		ViewCount:       int(viewCount),
		ClickCount:      int(clickCount),
		UniqueViewers:   int(uniqueViewers),
		UniqueClickers:  int(uniqueClickers),
		CTR:             ctr,
		Period:          period,
	}

	logger.Info("取得廣告統計成功",
		zap.Int("ad_id", adID),
		zap.String("period", period),
		zap.Int("views", stats.ViewCount),
		zap.Int("clicks", stats.ClickCount),
	)

	return stats, nil
}
