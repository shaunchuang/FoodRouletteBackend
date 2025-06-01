package domain

import (
	"time"
)

// GameType 遊戲類型
type GameType string

const (
	GameTypeRoulette GameType = "roulette" // 餐廳輪盤
	GameTypeDice     GameType = "dice"     // 骰子決定法
	GameTypeTarot    GameType = "tarot"    // 塔羅占卜
	GameTypePuzzle   GameType = "puzzle"   // 拼圖配對
	GameTypeMap      GameType = "map"      // 地圖尋寶
)

// GameSession 遊戲會話
type GameSession struct {
	ID                 string                   `json:"id" db:"id"` // UUID
	UserID             int                      `json:"user_id" db:"user_id"`
	GameType           GameType                 `json:"game_type" db:"game_type"`
	Status             string                   `json:"status" db:"status"`                             // pending, playing, completed
	ResultRestaurantID *int                     `json:"result_restaurant_id" db:"result_restaurant_id"` // 結果餐廳 ID
	Result             *RestaurantWithDistance  `json:"result"`                                         // 遊戲結果
	Restaurants        []RestaurantWithDistance `json:"restaurants"`                                    // 參與遊戲的餐廳列表
	Advertisements     []Advertisement          `json:"advertisements"`                                 // 顯示的廣告
	StartedAt          time.Time                `json:"started_at" db:"started_at"`
	CompletedAt        *time.Time               `json:"completed_at" db:"completed_at"`
	CreatedAt          time.Time                `json:"created_at" db:"created_at"`
}

// StartGameRequest 開始遊戲請求
type StartGameRequest struct {
	GameType  GameType `json:"game_type" validate:"required"`
	Latitude  float64  `json:"latitude" validate:"required,latitude"`
	Longitude float64  `json:"longitude" validate:"required,longitude"`
	Radius    int      `json:"radius" validate:"min=100,max=10000"` // 搜尋半徑（公尺）
}

// GameResult 遊戲結果
type GameResult struct {
	SessionID          string                  `json:"session_id"`
	SelectedRestaurant *RestaurantWithDistance `json:"selected_restaurant"`
	ClickedAd          *Advertisement          `json:"clicked_ad,omitempty"` // 如果點擊了廣告
	CompletedAt        time.Time               `json:"completed_at"`
}

// CompleteGameRequest 完成遊戲請求
type CompleteGameRequest struct {
	SessionID            string `json:"session_id" validate:"required"`
	SelectedRestaurantID int    `json:"selected_restaurant_id" validate:"required"`
	ClickedAdID          *int   `json:"clicked_ad_id,omitempty"` // 如果有點擊廣告
}
