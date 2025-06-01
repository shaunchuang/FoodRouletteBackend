package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTService JWT 認證服務
type JWTService struct {
	secretKey []byte
}

// NewJWTService 建立 JWT 服務
func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
	}
}

// Claims JWT Claims 結構
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// HashPassword 加密密碼
func (s *JWTService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 驗證密碼
func (s *JWTService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateToken 生成 JWT Token
func (s *JWTService) GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24小時過期

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "food-roulette-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 驗證 JWT Token
func (s *JWTService) ValidateToken(tokenString string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無效的簽名方法")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("無效的 token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("token 已過期")
	}

	return claims.UserID, nil
}
