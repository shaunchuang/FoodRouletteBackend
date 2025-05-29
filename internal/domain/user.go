package domain

import (
	"time"
)

// User 使用者實體
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Username  string    `json:"username" db:"username" validate:"required,min=3,max=50"`
	Password  string    `json:"-" db:"password"` // 不在 JSON 中顯示密碼
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserLocation 使用者位置資訊
type UserLocation struct {
	UserID    int     `json:"user_id" db:"user_id"`
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required,longitude"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest 建立使用者請求
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest 登入請求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateLocationRequest 更新位置請求
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}