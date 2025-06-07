package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
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
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
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

// GenerateToken 生成 JWT Token（實作介面）
func (s *JWTService) GenerateToken(userID int, email, username, role, provider string) (string, error) {
	return s.GenerateTokenWithDetails(userID, email, username, role, provider)
}

// GenerateTokenLegacy 生成 JWT Token（向後相容版本）
func (s *JWTService) GenerateTokenLegacy(userID int) (string, error) {
	return s.GenerateTokenWithDetails(userID, "", "", "user", "local")
}

// GenerateTokenWithDetails 生成詳細的 JWT Token
func (s *JWTService) GenerateTokenWithDetails(userID int, email, username, role, provider string) (string, error) {
	expirationTime := time.Now().Add(168 * time.Hour) // 7天過期

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		Provider: provider,
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
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}

// GoogleUserInfo Google 使用者資訊結構
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// VerifyGoogleToken 驗證 Google ID Token
func (s *JWTService) VerifyGoogleToken(idToken string) (*domain.OAuthUserInfo, error) {
	ctx := context.Background()

	// 建立 OAuth2 服務
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, errors.New("failed to create OAuth2 service")
	}

	// 驗證 token
	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(idToken).Do()
	if err != nil {
		return nil, errors.New("invalid Google ID token")
	}

	// 檢查 token 是否過期
	if tokenInfo.ExpiresIn <= 0 {
		return nil, errors.New("Google ID token has expired")
	}

	// 使用 access token 獲取詳細使用者資訊
	userInfoURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", idToken)
	resp, err := http.Get(userInfoURL)
	if err != nil {
		// 如果無法獲取詳細資訊，就使用基本資訊
		return &domain.OAuthUserInfo{
			Email:    tokenInfo.Email,
			Name:     tokenInfo.Email, // 使用 email 作為 fallback
			Picture:  "",
			Subject:  tokenInfo.UserId,
			Provider: "google",
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 如果無法獲取詳細資訊，就使用基本資訊
		return &domain.OAuthUserInfo{
			Email:    tokenInfo.Email,
			Name:     tokenInfo.Email, // 使用 email 作為 fallback
			Picture:  "",
			Subject:  tokenInfo.UserId,
			Provider: "google",
		}, nil
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		// 如果解析失敗，就使用基本資訊
		return &domain.OAuthUserInfo{
			Email:    tokenInfo.Email,
			Name:     tokenInfo.Email, // 使用 email 作為 fallback
			Picture:  "",
			Subject:  tokenInfo.UserId,
			Provider: "google",
		}, nil
	}

	return &domain.OAuthUserInfo{
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Picture:  userInfo.Picture,
		Subject:  userInfo.ID,
		Provider: "google",
	}, nil
}

// AppleTokenClaims Apple ID Token Claims 結構
type AppleTokenClaims struct {
	Issuer        string `json:"iss"`
	Subject       string `json:"sub"`
	Audience      string `json:"aud"`
	IssuedAt      int64  `json:"iat"`
	ExpiresAt     int64  `json:"exp"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	jwt.RegisteredClaims
}

// VerifyAppleToken 驗證 Apple ID Token
func (s *JWTService) VerifyAppleToken(idToken string) (*domain.OAuthUserInfo, error) {
	// 解析 JWT token 但不驗證簽名（簡化版實作）
	// 在生產環境中，應該從 Apple 獲取公鑰並驗證簽名
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, &AppleTokenClaims{})
	if err != nil {
		return nil, errors.New("invalid Apple ID token format")
	}

	claims, ok := token.Claims.(*AppleTokenClaims)
	if !ok {
		return nil, errors.New("invalid Apple ID token claims")
	}

	// 檢查 token 是否過期
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("Apple ID token has expired")
	}

	// 檢查 issuer
	if claims.Issuer != "https://appleid.apple.com" {
		return nil, errors.New("invalid Apple ID token issuer")
	}

	// Apple 通常不提供 name，使用 email 作為替代
	name := claims.Email
	if name == "" {
		name = "Apple User"
	}

	return &domain.OAuthUserInfo{
		Email:    claims.Email,
		Name:     name,
		Picture:  "", // Apple 通常不提供頭像
		Subject:  claims.Subject,
		Provider: "apple",
	}, nil
}
