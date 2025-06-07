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

// UserRepository PostgreSQL 使用者資料庫操作實作
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 建立使用者 Repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create 建立新使用者
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, username, password, provider, role, status, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		user.Password,
		user.Provider,
		user.Role,
		user.Status,
		user.EmailVerified,
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		logger.Error("建立使用者失敗", zap.Error(err), zap.String("email", user.Email))
		return err
	}

	user.CreatedAt = now
	user.UpdatedAt = now

	logger.Info("使用者建立成功", zap.Int("user_id", user.ID), zap.String("email", user.Email))
	return nil
}

// GetByID 根據 ID 取得使用者
func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, provider, provider_id, avatar, 
		       role, status, email_verified, last_login_at, failed_login_attempts, 
		       locked_until, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Provider,
		&user.ProviderID,
		&user.Avatar,
		&user.Role,
		&user.Status,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("使用者不存在")
		}
		logger.Error("取得使用者失敗", zap.Error(err), zap.Int("user_id", id))
		return nil, err
	}

	return user, nil
}

// GetByEmail 根據 Email 取得使用者
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, provider, provider_id, avatar, 
		       role, status, email_verified, last_login_at, failed_login_attempts, 
		       locked_until, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Provider,
		&user.ProviderID,
		&user.Avatar,
		&user.Role,
		&user.Status,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("使用者不存在")
		}
		logger.Error("根據 Email 取得使用者失敗", zap.Error(err), zap.String("email", email))
		return nil, err
	}

	return user, nil
}

// Update 更新使用者資訊
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2, password = $3, updated_at = $4
		WHERE id = $5`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Username,
		user.Password,
		now,
		user.ID,
	)

	if err != nil {
		logger.Error("更新使用者失敗", zap.Error(err), zap.Int("user_id", user.ID))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("使用者不存在")
	}

	user.UpdatedAt = now
	logger.Info("使用者更新成功", zap.Int("user_id", user.ID))
	return nil
}

// UpdateLocation 更新使用者位置
func (r *UserRepository) UpdateLocation(ctx context.Context, userID int, location *domain.UserLocation) error {
	query := `
		INSERT INTO user_locations (user_id, latitude, longitude, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			updated_at = EXCLUDED.updated_at`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		userID,
		location.Latitude,
		location.Longitude,
		now,
	)

	if err != nil {
		logger.Error("更新使用者位置失敗", zap.Error(err), zap.Int("user_id", userID))
		return err
	}

	location.UserID = userID
	location.UpdatedAt = now

	logger.Info("使用者位置更新成功", zap.Int("user_id", userID))
	return nil
}

// GetLocation 取得使用者位置
func (r *UserRepository) GetLocation(ctx context.Context, userID int) (*domain.UserLocation, error) {
	query := `
		SELECT user_id, latitude, longitude, updated_at
		FROM user_locations
		WHERE user_id = $1`

	location := &domain.UserLocation{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&location.UserID,
		&location.Latitude,
		&location.Longitude,
		&location.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("使用者位置未找到")
		}
		logger.Error("取得使用者位置失敗", zap.Error(err))
		return nil, err
	}

	return location, nil
}

// GetByProviderID 根據 Provider 和 Provider ID 取得使用者
func (r *UserRepository) GetByProviderID(ctx context.Context, provider, providerID string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, provider, provider_id, avatar, 
		       role, status, email_verified, last_login_at, failed_login_attempts, 
		       locked_until, created_at, updated_at
		FROM users
		WHERE provider = $1 AND provider_id = $2`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Provider,
		&user.ProviderID,
		&user.Avatar,
		&user.Role,
		&user.Status,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("使用者未找到")
		}
		logger.Error("根據 Provider 查找使用者失敗", zap.Error(err))
		return nil, err
	}

	return user, nil
}

// UpdateProviderInfo 更新使用者的第三方登入資訊
func (r *UserRepository) UpdateProviderInfo(ctx context.Context, userID int, provider, providerID string) error {
	query := `
		UPDATE users
		SET provider = $1, provider_id = $2, updated_at = $3
		WHERE id = $4`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, provider, providerID, now, userID)
	if err != nil {
		logger.Error("更新使用者 Provider 資訊失敗", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("檢查更新結果失敗", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return errors.New("使用者未找到")
	}

	logger.Info("使用者 Provider 資訊更新成功", zap.Int("user_id", userID))
	return nil
}

// CreateOAuthUser 建立 OAuth 使用者
func (r *UserRepository) CreateOAuthUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, username, provider, provider_id, avatar, role, status, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		user.Provider,
		user.ProviderID,
		user.Avatar,
		user.Role,
		user.Status,
		user.EmailVerified,
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		logger.Error("建立 OAuth 使用者失敗", zap.Error(err))
		return err
	}

	user.CreatedAt = now
	user.UpdatedAt = now

	logger.Info("OAuth 使用者建立成功", zap.Int("user_id", user.ID), zap.String("email", user.Email))
	return nil
}
