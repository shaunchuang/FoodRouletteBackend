package usecase

import (
	"context"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
)

// UserRepository 使用者資料庫操作介面
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateLocation(ctx context.Context, userID int, location *domain.UserLocation) error
	GetLocation(ctx context.Context, userID int) (*domain.UserLocation, error)
}

// RestaurantRepository 餐廳資料庫操作介面
type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *domain.Restaurant) error
	GetByID(ctx context.Context, id int) (*domain.Restaurant, error)
	SearchNearby(ctx context.Context, params *domain.RestaurantSearchParams) ([]domain.RestaurantWithDistance, error)
	Update(ctx context.Context, restaurant *domain.Restaurant) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Restaurant, error)
}

// FavoriteRepository 最愛餐廳資料庫操作介面
type FavoriteRepository interface {
	Add(ctx context.Context, userID int, request *domain.AddFavoriteRequest) error
	Remove(ctx context.Context, userID, restaurantID int) error
	GetByUserID(ctx context.Context, userID int) ([]domain.FavoriteRestaurant, error)
	IsExists(ctx context.Context, userID, restaurantID int) (bool, error)
}

// GameRepository 遊戲會話資料庫操作介面
type GameRepository interface {
	CreateSession(ctx context.Context, session *domain.GameSession) error
	GetSessionByID(ctx context.Context, sessionID string) (*domain.GameSession, error)
	UpdateSession(ctx context.Context, session *domain.GameSession) error
	GetUserSessions(ctx context.Context, userID int, limit, offset int) ([]domain.GameSession, error)
}

// AdvertisementRepository 廣告資料庫操作介面
type AdvertisementRepository interface {
	Create(ctx context.Context, ad *domain.Advertisement) error
	GetByID(ctx context.Context, id int) (*domain.Advertisement, error)
	GetActiveAds(ctx context.Context, limit int) ([]domain.Advertisement, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Advertisement, error)
	Update(ctx context.Context, ad *domain.Advertisement) error
	RecordView(ctx context.Context, view *domain.AdView) error
	RecordClick(ctx context.Context, click *domain.AdClick) error
	GetStatistics(ctx context.Context, adID int, period string) (*domain.AdStatistics, error)
}

// ExternalAPIService 外部 API 服務介面
type ExternalAPIService interface {
	SearchNearbyRestaurants(ctx context.Context, lat, lng float64, radius int) ([]domain.Restaurant, error)
	GetRestaurantDetails(ctx context.Context, googleID string) (*domain.Restaurant, error)
}

// AuthService 認證服務介面
type AuthService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) bool
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (int, error)
}

// UserService 使用者服務介面
type UserService interface {
	Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userID int) (*domain.User, error)
	UpdateLocation(ctx context.Context, userID int, req *domain.UpdateLocationRequest) error
	GetLocation(ctx context.Context, userID int) (*domain.UserLocation, error)
}
