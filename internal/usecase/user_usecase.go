package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// UserUseCase 使用者業務邏輯
type UserUseCase struct {
	userRepo UserRepository
	authSvc  AuthService
}

// NewUserUseCase 建立使用者用例
func NewUserUseCase(userRepo UserRepository, authSvc AuthService) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

// Register 使用者註冊
func (uc *UserUseCase) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// 檢查信箱是否已存在
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("電子郵件已被使用")
	}

	// 加密密碼
	hashedPassword, err := uc.authSvc.HashPassword(req.Password)
	if err != nil {
		logger.Error("密碼加密失敗", zap.Error(err))
		return nil, errors.New("密碼處理失敗")
	}

	// 建立使用者
	user := &domain.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.Error("建立使用者失敗", zap.Error(err))
		return nil, errors.New("建立使用者失敗")
	}

	logger.Info("使用者註冊成功", zap.String("email", user.Email))
	return user, nil
}

// Login 使用者登入
func (uc *UserUseCase) Login(ctx context.Context, req *domain.LoginRequest) (string, error) {
	// 根據信箱查找使用者
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("查找使用者失敗", zap.Error(err))
		return "", errors.New("電子郵件或密碼錯誤")
	}

	// 驗證密碼
	if !uc.authSvc.VerifyPassword(user.Password, req.Password) {
		logger.Warn("密碼驗證失敗", zap.String("email", req.Email))
		return "", errors.New("電子郵件或密碼錯誤")
	}

	// 生成 JWT Token
	token, err := uc.authSvc.GenerateToken(user.ID)
	if err != nil {
		logger.Error("生成 Token 失敗", zap.Error(err))
		return "", errors.New("登入失敗")
	}

	logger.Info("使用者登入成功", zap.String("email", user.Email))
	return token, nil
}

// GetProfile 取得使用者資料
func (uc *UserUseCase) GetProfile(ctx context.Context, userID int) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Error("取得使用者資料失敗", zap.Error(err), zap.Int("user_id", userID))
		return nil, errors.New("使用者不存在")
	}

	return user, nil
}

// UpdateLocation 更新使用者位置
func (uc *UserUseCase) UpdateLocation(ctx context.Context, userID int, req *domain.UpdateLocationRequest) error {
	location := &domain.UserLocation{
		UserID:    userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.UpdateLocation(ctx, userID, location); err != nil {
		logger.Error("更新使用者位置失敗", zap.Error(err), zap.Int("user_id", userID))
		return errors.New("更新位置失敗")
	}

	logger.Info("使用者位置更新成功", zap.Int("user_id", userID))
	return nil
}

// GetLocation 取得使用者位置
func (uc *UserUseCase) GetLocation(ctx context.Context, userID int) (*domain.UserLocation, error) {
	location, err := uc.userRepo.GetLocation(ctx, userID)
	if err != nil {
		logger.Error("取得使用者位置失敗", zap.Error(err), zap.Int("user_id", userID))
		return nil, errors.New("取得位置失敗")
	}

	return location, nil
}