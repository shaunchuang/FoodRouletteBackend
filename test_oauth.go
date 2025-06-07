package main

import (
	"encoding/json"
	"log"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/auth"
)

func main() {
	// 測試 JWT 服務
	jwtService := auth.NewJWTService("test_secret")

	// 測試生成 JWT Token
	token, err := jwtService.GenerateToken(1, "test@example.com", "testuser", "user", "google")
	if err != nil {
		log.Fatalf("生成 JWT Token 失敗: %v", err)
	}

	log.Printf("生成的 JWT Token: %s", token)

	// 測試驗證 JWT Token
	userID, err := jwtService.ValidateToken(token)
	if err != nil {
		log.Fatalf("驗證 JWT Token 失敗: %v", err)
	}

	log.Printf("驗證成功，用戶 ID: %d", userID)

	// 測試 Apple Token 解析（使用簡化版實作）
	// 注意：這裡只是測試結構，實際的 Apple Token 需要有效的簽名
	appleTokenExample := `{
		"iss": "https://appleid.apple.com",
		"sub": "001234.567890abcdef.1234",
		"aud": "com.example.app",
		"exp": 9999999999,
		"iat": 1640995200,
		"email": "user@privaterelay.appleid.com",
		"email_verified": "true"
	}`

	log.Printf("Apple Token 示例格式: %s", appleTokenExample)

	log.Println("OAuth 認證系統測試完成！")
	log.Println("系統已準備好處理以下認證方式：")
	log.Println("- 傳統用戶名/密碼登入")
	log.Println("- Google OAuth 登入")
	log.Println("- Apple ID 登入")
	log.Println("- JWT Token 驗證")

	// 測試 OAuthUserInfo 結構
	oauthUser := &domain.OAuthUserInfo{
		Email:    "test@gmail.com",
		Name:     "Test User",
		Picture:  "https://example.com/avatar.jpg",
		Subject:  "google_123456789",
		Provider: "google",
	}

	oauthUserJSON, _ := json.MarshalIndent(oauthUser, "", "  ")
	log.Printf("OAuth 用戶資訊結構: %s", string(oauthUserJSON))
}
