package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// RestaurantRepository PostgreSQL 餐廳資料庫操作實作
type RestaurantRepository struct {
	db *sql.DB
}

// NewRestaurantRepository 建立餐廳 Repository
func NewRestaurantRepository(db *sql.DB) *RestaurantRepository {
	return &RestaurantRepository{
		db: db,
	}
}

// Create 建立新餐廳
func (r *RestaurantRepository) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	query := `
		INSERT INTO restaurants (name, address, latitude, longitude, phone, rating, price_level, cuisine, is_active, google_id, image_url, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		restaurant.Name,
		restaurant.Address,
		restaurant.Latitude,
		restaurant.Longitude,
		restaurant.Phone,
		restaurant.Rating,
		restaurant.PriceLevel,
		restaurant.Cuisine,
		restaurant.IsActive,
		restaurant.GoogleID,
		restaurant.ImageURL,
		restaurant.Description,
		now,
		now,
	).Scan(&restaurant.ID)

	if err != nil {
		logger.Error("建立餐廳失敗", zap.Error(err), zap.String("name", restaurant.Name))
		return err
	}

	restaurant.CreatedAt = now
	restaurant.UpdatedAt = now

	logger.Info("餐廳建立成功", zap.Int("restaurant_id", restaurant.ID), zap.String("name", restaurant.Name))
	return nil
}

// GetByID 根據 ID 取得餐廳
func (r *RestaurantRepository) GetByID(ctx context.Context, id int) (*domain.Restaurant, error) {
	query := `
		SELECT id, name, address, latitude, longitude, phone, rating, price_level, cuisine, is_active, google_id, image_url, description, created_at, updated_at
		FROM restaurants
		WHERE id = $1 AND is_active = TRUE`

	restaurant := &domain.Restaurant{}
	var phone, googleID, imageURL, description sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Address,
		&restaurant.Latitude,
		&restaurant.Longitude,
		&phone,
		&restaurant.Rating,
		&restaurant.PriceLevel,
		&restaurant.Cuisine,
		&restaurant.IsActive,
		&googleID,
		&imageURL,
		&description,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("餐廳不存在")
		}
		logger.Error("取得餐廳失敗", zap.Error(err), zap.Int("restaurant_id", id))
		return nil, err
	}

	// 處理可為空的欄位
	if phone.Valid {
		restaurant.Phone = phone.String
	}
	if googleID.Valid {
		restaurant.GoogleID = googleID.String
	}
	if imageURL.Valid {
		restaurant.ImageURL = imageURL.String
	}
	if description.Valid {
		restaurant.Description = description.String
	}

	return restaurant, nil
}

// SearchNearby 搜尋附近餐廳
func (r *RestaurantRepository) SearchNearby(ctx context.Context, params *domain.RestaurantSearchParams) ([]domain.RestaurantWithDistance, error) {
	// 建立基本查詢
	baseQuery := `
		SELECT id, name, address, latitude, longitude, phone, rating, price_level, cuisine, is_active, google_id, image_url, description,
		       (6371 * acos(cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))) AS distance
		FROM restaurants
		WHERE is_active = TRUE`

	args := []interface{}{params.Latitude, params.Longitude}
	argIndex := 2

	// 添加距離篩選
	baseQuery += fmt.Sprintf(" AND (6371 * acos(cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))) <= $%d", argIndex+1)
	args = append(args, float64(params.Radius)/1000) // 轉換為公里
	argIndex++

	// 添加料理類型篩選
	if params.Cuisine != "" {
		baseQuery += fmt.Sprintf(" AND cuisine ILIKE $%d", argIndex+1)
		args = append(args, "%"+params.Cuisine+"%")
		argIndex++
	}

	// 添加評分篩選
	if params.MinRating > 0 {
		baseQuery += fmt.Sprintf(" AND rating >= $%d", argIndex+1)
		args = append(args, params.MinRating)
		argIndex++
	}

	// 排序和限制
	baseQuery += " ORDER BY distance"
	if params.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex+1)
		args = append(args, params.Limit)
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		logger.Error("搜尋附近餐廳失敗", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var restaurants []domain.RestaurantWithDistance
	for rows.Next() {
		var restaurant domain.RestaurantWithDistance
		var phone, googleID, imageURL, description sql.NullString

		err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.Latitude,
			&restaurant.Longitude,
			&phone,
			&restaurant.Rating,
			&restaurant.PriceLevel,
			&restaurant.Cuisine,
			&restaurant.IsActive,
			&googleID,
			&imageURL,
			&description,
			&restaurant.Distance,
		)

		if err != nil {
			logger.Error("掃描餐廳資料失敗", zap.Error(err))
			continue
		}

		// 處理可為空的欄位
		if phone.Valid {
			restaurant.Phone = phone.String
		}
		if googleID.Valid {
			restaurant.GoogleID = googleID.String
		}
		if imageURL.Valid {
			restaurant.ImageURL = imageURL.String
		}
		if description.Valid {
			restaurant.Description = description.String
		}

		// 轉換距離為公尺
		restaurant.Distance = restaurant.Distance * 1000

		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		logger.Error("處理餐廳查詢結果失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("搜尋附近餐廳成功",
		zap.Float64("latitude", params.Latitude),
		zap.Float64("longitude", params.Longitude),
		zap.Int("radius", params.Radius),
		zap.Int("count", len(restaurants)),
	)

	return restaurants, nil
}

// Update 更新餐廳資訊
func (r *RestaurantRepository) Update(ctx context.Context, restaurant *domain.Restaurant) error {
	query := `
		UPDATE restaurants
		SET name = $1, address = $2, latitude = $3, longitude = $4, phone = $5, rating = $6, 
		    price_level = $7, cuisine = $8, is_active = $9, google_id = $10, image_url = $11, 
		    description = $12, updated_at = $13
		WHERE id = $14`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		restaurant.Name,
		restaurant.Address,
		restaurant.Latitude,
		restaurant.Longitude,
		restaurant.Phone,
		restaurant.Rating,
		restaurant.PriceLevel,
		restaurant.Cuisine,
		restaurant.IsActive,
		restaurant.GoogleID,
		restaurant.ImageURL,
		restaurant.Description,
		now,
		restaurant.ID,
	)

	if err != nil {
		logger.Error("更新餐廳失敗", zap.Error(err), zap.Int("restaurant_id", restaurant.ID))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("餐廳不存在")
	}

	restaurant.UpdatedAt = now
	logger.Info("餐廳更新成功", zap.Int("restaurant_id", restaurant.ID))
	return nil
}

// GetAll 取得所有餐廳（管理功能）
func (r *RestaurantRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Restaurant, error) {
	query := `
		SELECT id, name, address, latitude, longitude, phone, rating, price_level, cuisine, is_active, google_id, image_url, description, created_at, updated_at
		FROM restaurants
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		logger.Error("取得餐廳清單失敗", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var restaurants []domain.Restaurant
	for rows.Next() {
		var restaurant domain.Restaurant
		var phone, googleID, imageURL, description sql.NullString

		err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.Latitude,
			&restaurant.Longitude,
			&phone,
			&restaurant.Rating,
			&restaurant.PriceLevel,
			&restaurant.Cuisine,
			&restaurant.IsActive,
			&googleID,
			&imageURL,
			&description,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)

		if err != nil {
			logger.Error("掃描餐廳資料失敗", zap.Error(err))
			continue
		}

		// 處理可為空的欄位
		if phone.Valid {
			restaurant.Phone = phone.String
		}
		if googleID.Valid {
			restaurant.GoogleID = googleID.String
		}
		if imageURL.Valid {
			restaurant.ImageURL = imageURL.String
		}
		if description.Valid {
			restaurant.Description = description.String
		}

		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		logger.Error("處理餐廳查詢結果失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("取得餐廳清單成功", zap.Int("count", len(restaurants)))
	return restaurants, nil
}

// calculateDistance 計算兩點之間的距離（使用 Haversine 公式）
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // 地球半徑（公里）

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c * 1000 // 轉換為公尺
}
