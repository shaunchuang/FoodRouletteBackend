package domain

import (
	"time"
)

// Advertisement 廣告實體
type Advertisement struct {
	ID           int       `json:"id" db:"id"`
	RestaurantID int       `json:"restaurant_id" db:"restaurant_id"`
	Title        string    `json:"title" db:"title" validate:"required,max=100"`
	Content      string    `json:"content" db:"content" validate:"required,max=500"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	TargetURL    string    `json:"target_url" db:"target_url"` // 點擊後跳轉的 URL
	IsActive     bool      `json:"is_active" db:"is_active"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	Priority     int       `json:"priority" db:"priority"` // 廣告優先級，數字越大優先級越高
	ClickCount   int       `json:"click_count" db:"click_count"`
	ViewCount    int       `json:"view_count" db:"view_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// AdClick 廣告點擊記錄
type AdClick struct {
	ID              int       `json:"id" db:"id"`
	AdvertisementID int       `json:"advertisement_id" db:"advertisement_id"`
	UserID          int       `json:"user_id" db:"user_id"`
	GameSessionID   string    `json:"game_session_id" db:"game_session_id"`
	ClickedAt       time.Time `json:"clicked_at" db:"clicked_at"`
	IPAddress       string    `json:"ip_address" db:"ip_address"`
	UserAgent       string    `json:"user_agent" db:"user_agent"`
}

// AdView 廣告瀏覽記錄
type AdView struct {
	ID              int       `json:"id" db:"id"`
	AdvertisementID int       `json:"advertisement_id" db:"advertisement_id"`
	UserID          int       `json:"user_id" db:"user_id"`
	GameSessionID   string    `json:"game_session_id" db:"game_session_id"`
	ViewedAt        time.Time `json:"viewed_at" db:"viewed_at"`
	IPAddress       string    `json:"ip_address" db:"ip_address"`
	UserAgent       string    `json:"user_agent" db:"user_agent"`
}

// CreateAdRequest 建立廣告請求
type CreateAdRequest struct {
	RestaurantID int       `json:"restaurant_id" validate:"required"`
	Title        string    `json:"title" validate:"required,max=100"`
	Content      string    `json:"content" validate:"required,max=500"`
	ImageURL     string    `json:"image_url" validate:"url"`
	TargetURL    string    `json:"target_url" validate:"url"`
	StartDate    time.Time `json:"start_date" validate:"required"`
	EndDate      time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	Priority     int       `json:"priority" validate:"min=1,max=10"`
}

// AdStatistics 廣告統計資訊
type AdStatistics struct {
	AdvertisementID int     `json:"advertisement_id"`
	ViewCount       int     `json:"view_count"`
	ClickCount      int     `json:"click_count"`
	UniqueViewers   int     `json:"unique_viewers"`  // 獨立瀏覽用戶數
	UniqueClickers  int     `json:"unique_clickers"` // 獨立點擊用戶數
	CTR             float64 `json:"ctr"`             // Click Through Rate
	Period          string  `json:"period"`          // daily, weekly, monthly
}
