package domain

import (
	"time"
)

// Restaurant 餐廳實體
type Restaurant struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name" validate:"required,max=100"`
	Address     string  `json:"address" db:"address" validate:"required,max=255"`
	Latitude    float64 `json:"latitude" db:"latitude" validate:"required,latitude"`
	Longitude   float64 `json:"longitude" db:"longitude" validate:"required,longitude"`
	Phone       string  `json:"phone" db:"phone"`
	Rating      float32 `json:"rating" db:"rating" validate:"min=0,max=5"`
	PriceLevel  int     `json:"price_level" db:"price_level" validate:"min=1,max=4"` // 1-4 價位等級
	Cuisine     string  `json:"cuisine" db:"cuisine"`                                 // 料理類型
	IsActive    bool    `json:"is_active" db:"is_active"`
	GoogleID    string  `json:"google_id" db:"google_id"`    // Google Places ID
	ImageURL    string  `json:"image_url" db:"image_url"`    // 餐廳圖片
	Description string  `json:"description" db:"description"` // 餐廳描述
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// FavoriteRestaurant 使用者最愛餐廳
type FavoriteRestaurant struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	RestaurantID int       `json:"restaurant_id" db:"restaurant_id"`
	Notes        string    `json:"notes" db:"notes"` // 使用者備註
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// RestaurantSearchParams 餐廳搜尋參數
type RestaurantSearchParams struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
	Radius    int     `json:"radius" validate:"min=100,max=10000"` // 搜尋半徑（公尺）
	Cuisine   string  `json:"cuisine"`                             // 料理類型篩選
	MinRating float32 `json:"min_rating" validate:"min=0,max=5"`   // 最低評分
	Limit     int     `json:"limit" validate:"min=1,max=50"`       // 結果數量限制
}

// AddFavoriteRequest 新增最愛餐廳請求
type AddFavoriteRequest struct {
	RestaurantID int    `json:"restaurant_id" validate:"required"`
	Notes        string `json:"notes" validate:"max=500"`
}

// RestaurantWithDistance 包含距離資訊的餐廳
type RestaurantWithDistance struct {
	Restaurant
	Distance float64 `json:"distance"` // 距離（公尺）
}