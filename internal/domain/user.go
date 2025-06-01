package domain

import (
	"time"
)

// UserRole 使用者角色類型
type UserRole string

const (
	UserRoleUser      UserRole = "user"
	UserRoleAdmin     UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
)

// UserStatus 使用者狀態類型
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// User 使用者實體
type User struct {
	ID                  int        `json:"id" db:"id"`
	Email               string     `json:"email" db:"email" validate:"required,email"`
	Username            string     `json:"username" db:"username" validate:"required,min=3,max=50"`
	Password            string     `json:"-" db:"password"` // 不在 JSON 中顯示密碼
	Role                UserRole   `json:"role" db:"role"`
	Status              UserStatus `json:"status" db:"status"`
	EmailVerified       bool       `json:"email_verified" db:"email_verified"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	FailedLoginAttempts int        `json:"-" db:"failed_login_attempts"` // 不在 JSON 中顯示
	LockedUntil         *time.Time `json:"-" db:"locked_until"`          // 不在 JSON 中顯示
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// IsAdmin 檢查使用者是否為管理員
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsModerator 檢查使用者是否為版主
func (u *User) IsModerator() bool {
	return u.Role == UserRoleModerator
}

// HasAdminAccess 檢查使用者是否有管理員存取權限
func (u *User) HasAdminAccess() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleModerator
}

// IsActive 檢查使用者是否為活躍狀態
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsLocked 檢查使用者是否被鎖定
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// UserLocation 使用者位置資訊
type UserLocation struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Latitude  float64   `json:"latitude" db:"latitude" validate:"required,latitude"`
	Longitude float64   `json:"longitude" db:"longitude" validate:"required,longitude"`
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
