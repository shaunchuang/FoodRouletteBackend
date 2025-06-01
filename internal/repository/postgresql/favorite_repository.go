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

// FavoriteRepository PostgreSQL 最愛餐廳資料庫操作實作
type FavoriteRepository struct {
	db *sql.DB
}

// NewFavoriteRepository 建立最愛餐廳 Repository
func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
	return &FavoriteRepository{
		db: db,
	}
}

// Add 新增最愛餐廳
func (r *FavoriteRepository) Add(ctx context.Context, userID int, request *domain.AddFavoriteRequest) error {
	query := `
		INSERT INTO favorite_restaurants (user_id, restaurant_id, notes, created_at)
		VALUES ($1, $2, $3, $4)`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		userID,
		request.RestaurantID,
		request.Notes,
		now,
	)

	if err != nil {
		logger.Error("新增最愛餐廳失敗", zap.Error(err), zap.Int("user_id", userID), zap.Int("restaurant_id", request.RestaurantID))
		return err
	}

	logger.Info("新增最愛餐廳成功", zap.Int("user_id", userID), zap.Int("restaurant_id", request.RestaurantID))
	return nil
}

// Remove 移除最愛餐廳
func (r *FavoriteRepository) Remove(ctx context.Context, userID, restaurantID int) error {
	query := `
		DELETE FROM favorite_restaurants
		WHERE user_id = $1 AND restaurant_id = $2`

	result, err := r.db.ExecContext(ctx, query, userID, restaurantID)
	if err != nil {
		logger.Error("移除最愛餐廳失敗", zap.Error(err), zap.Int("user_id", userID), zap.Int("restaurant_id", restaurantID))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("最愛餐廳不存在")
	}

	logger.Info("移除最愛餐廳成功", zap.Int("user_id", userID), zap.Int("restaurant_id", restaurantID))
	return nil
}

// GetByUserID 取得使用者的最愛餐廳清單
func (r *FavoriteRepository) GetByUserID(ctx context.Context, userID int) ([]domain.FavoriteRestaurant, error) {
	query := `
		SELECT fr.id, fr.user_id, fr.restaurant_id, fr.notes, fr.created_at
		FROM favorite_restaurants fr
		JOIN restaurants res ON fr.restaurant_id = res.id
		WHERE fr.user_id = $1 AND res.is_active = TRUE
		ORDER BY fr.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Error("取得使用者最愛餐廳失敗", zap.Error(err), zap.Int("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var favorites []domain.FavoriteRestaurant
	for rows.Next() {
		var favorite domain.FavoriteRestaurant
		var notes sql.NullString

		err := rows.Scan(
			&favorite.ID,
			&favorite.UserID,
			&favorite.RestaurantID,
			&notes,
			&favorite.CreatedAt,
		)

		if err != nil {
			logger.Error("掃描最愛餐廳資料失敗", zap.Error(err))
			continue
		}

		// 處理可為空的欄位
		if notes.Valid {
			favorite.Notes = notes.String
		}

		favorites = append(favorites, favorite)
	}

	if err = rows.Err(); err != nil {
		logger.Error("處理最愛餐廳查詢結果失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("取得使用者最愛餐廳成功", zap.Int("user_id", userID), zap.Int("count", len(favorites)))
	return favorites, nil
}

// IsExists 檢查是否已存在最愛餐廳
func (r *FavoriteRepository) IsExists(ctx context.Context, userID, restaurantID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM favorite_restaurants
			WHERE user_id = $1 AND restaurant_id = $2
		)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, restaurantID).Scan(&exists)
	if err != nil {
		logger.Error("檢查最愛餐廳是否存在失敗", zap.Error(err), zap.Int("user_id", userID), zap.Int("restaurant_id", restaurantID))
		return false, err
	}

	return exists, nil
}
