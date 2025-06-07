# 使用者驗證系統設計文件

## 概述

本文件規劃使用者登入系統，支援多平台（Android、iOS、Web）的第三方登入（Google、Apple ID），並在後端資料庫儲存對應使用者資料。

## 平台支援

- **Android App**: Google 登入
- **iOS App**: Apple ID 登入、Google 登入
- **Web 前端**: Google 登入、Apple ID 登入

## 資料庫設計

### Users 資料表欄位擴充

需要在現有的 [`users`](migrations/000001_init_schema.up.sql) 資料表中新增以下欄位：

```sql
-- 新增欄位
ALTER TABLE users ADD COLUMN provider VARCHAR(20) NULL;
ALTER TABLE users ADD COLUMN provider_id VARCHAR(255) NULL;
ALTER TABLE users ADD COLUMN avatar TEXT NULL;

-- 建立索引
CREATE INDEX idx_users_provider_id ON users(provider, provider_id);
CREATE UNIQUE INDEX idx_users_provider_unique ON users(provider, provider_id) WHERE provider IS NOT NULL;
```

#### 完整資料表結構

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255),  -- 第三方登入時可為 NULL
    provider VARCHAR(20),   -- 'google', 'apple', 'local'
    provider_id VARCHAR(255),  -- 第三方帳號唯一識別碼
    avatar TEXT,           -- 頭像 URL
    role VARCHAR(20) DEFAULT 'user',
    status VARCHAR(20) DEFAULT 'active',
    email_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API 設計

### 第三方登入 API

#### Google 登入
```http
POST /api/v1/auth/google
Content-Type: application/json

{
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Apple ID 登入
```http
POST /api/v1/auth/apple
Content-Type: application/json

{
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 傳統帳密登入
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### API 回應格式

#### 成功回應
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "username",
      "avatar": "https://example.com/avatar.jpg",
      "role": "user",
      "status": "active",
      "provider": "google"
    }
  }
}
```

#### 錯誤回應
```json
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "無效的 ID Token"
  }
}
```

## JWT 設計

### JWT Payload 結構
```json
{
  "sub": "1",           // User ID
  "email": "user@example.com",
  "username": "username",
  "role": "user",
  "provider": "google",
  "iat": 1640995200,    // 發行時間
  "exp": 1641601200     // 過期時間（建議 7 天）
}
```

### JWT 使用方式
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 登入流程

### 第三方登入流程圖

```mermaid
sequenceDiagram
    participant User as 使用者
    participant Client as 前端 App/Web
    participant Backend as 後端 API
    participant DB as 資料庫
    participant OAuth as Google/Apple

    User->>Client: 點擊 Google/Apple 登入
    Client->>OAuth: 啟動 OAuth 流程
    OAuth-->>Client: 回傳 id_token
    Client->>Backend: POST /api/v1/auth/{provider}
    note over Client,Backend: { id_token: "..." }
    
    Backend->>OAuth: 驗證 id_token
    OAuth-->>Backend: 回傳使用者資訊
    note over OAuth,Backend: { email, name, picture }
    
    Backend->>DB: 查詢 email 是否存在
    alt 使用者已存在
        Backend->>DB: 更新 provider, provider_id
        note over Backend,DB: 自動綁定第三方登入
    else 使用者不存在
        Backend->>DB: 建立新使用者
        note over Backend,DB: 設定 provider, provider_id
    end
    
    Backend->>Backend: 產生 JWT Token
    Backend-->>Client: 回傳 JWT + 使用者資料
    Client-->>User: 登入成功，跳轉主畫面
```

### 帳號綁定邏輯

1. **相同 Email 自動綁定**: 若資料庫已有相同 email 的使用者，系統會自動綁定第三方登入方式
2. **多重登入支援**: 使用者可同時使用密碼登入與第三方登入
3. **Provider 優先順序**: 
   - 首次註冊時的 provider 為主要 provider
   - 後續綁定的 provider 會更新但保留原有資訊

## 前端實作建議

### Android 實作
```kotlin
// Google Sign-In
private fun signInWithGoogle() {
    val signInIntent = googleSignInClient.signInIntent
    startActivityForResult(signInIntent, RC_SIGN_IN)
}

private fun handleSignInResult(task: Task<GoogleSignInAccount>) {
    try {
        val account = task.getResult(ApiException::class.java)
        val idToken = account?.idToken
        // 呼叫後端 API
        authRepository.googleLogin(idToken)
    } catch (e: ApiException) {
        // 處理錯誤
    }
}
```

### iOS 實作
```swift
// Google Sign-In
@objc func googleSignIn() {
    guard let presentingViewController = UIApplication.shared.windows.first?.rootViewController else { return }
    
    GIDSignIn.sharedInstance.signIn(withPresenting: presentingViewController) { [weak self] result, error in
        guard let user = result?.user,
              let idToken = user.idToken?.tokenString else { return }
        
        // 呼叫後端 API
        self?.authService.googleLogin(idToken: idToken)
    }
}

// Apple Sign-In
@objc func appleSignIn() {
    let request = ASAuthorizationAppleIDProvider().createRequest()
    request.requestedScopes = [.fullName, .email]
    
    let controller = ASAuthorizationController(authorizationRequests: [request])
    controller.delegate = self
    controller.presentationContextProvider = self
    controller.performRequests()
}
```

### Web 實作
```javascript
// Google Sign-In
function initializeGoogleSignIn() {
  gapi.load('auth2', () => {
    gapi.auth2.init({
      client_id: 'YOUR_GOOGLE_CLIENT_ID'
    });
  });
}

function signInWithGoogle() {
  const authInstance = gapi.auth2.getAuthInstance();
  authInstance.signIn().then(googleUser => {
    const idToken = googleUser.getAuthResponse().id_token;
    // 呼叫後端 API
    fetch('/api/v1/auth/google', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id_token: idToken })
    });
  });
}

// Apple Sign-In
function signInWithApple() {
  AppleID.auth.signIn({
    clientId: 'YOUR_APPLE_CLIENT_ID',
    redirectURI: 'YOUR_REDIRECT_URI',
    scope: 'name email'
  }).then(response => {
    // 呼叫後端 API
    fetch('/api/v1/auth/apple', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id_token: response.authorization.id_token })
    });
  });
}
```

## 後端實作重點

### Domain 層擴充

需要在 [`internal/domain/user.go`](internal/domain/user.go) 中新增：

```go
// OAuth 登入請求
type OAuthLoginRequest struct {
    IDToken string `json:"id_token" validate:"required"`
}

// OAuth 使用者資訊
type OAuthUserInfo struct {
    Email    string `json:"email"`
    Name     string `json:"name"`
    Picture  string `json:"picture"`
    Subject  string `json:"sub"` // provider_id
}
```

### Repository 層擴充

在 [`internal/repository/postgresql/user_repository.go`](internal/repository/postgresql/user_repository.go) 中新增：

```go
// GetByProviderID 根據 Provider 和 Provider ID 取得使用者
func (r *UserRepository) GetByProviderID(ctx context.Context, provider, providerID string) (*domain.User, error)

// UpdateProviderInfo 更新使用者的第三方登入資訊
func (r *UserRepository) UpdateProviderInfo(ctx context.Context, userID int, provider, providerID string) error
```

### Handler 層新增

需要新增 OAuth 相關的 handler：

```go
// GoogleLogin Google 登入
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
    var req domain.OAuthLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 驗證 Google ID Token
    userInfo, err := h.authUsecase.VerifyGoogleToken(req.IDToken)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid token"})
        return
    }
    
    // 處理使用者登入/註冊
    user, token, err := h.authUsecase.OAuthLogin(userInfo, "google")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "token": token,
        "user":  user,
    })
}
```

## 安全性考量

### Token 驗證
- **Google ID Token**: 使用 Google 公鑰驗證簽章
- **Apple ID Token**: 使用 Apple 公鑰驗證簽章
- **JWT Token**: 使用 HMAC-SHA256 或 RS256 演算法

### 安全性措施
1. **HTTPS 必要**: 所有 API 必須使用 HTTPS
2. **Token 過期**: JWT 設定適當過期時間（建議 7 天）
3. **Rate Limiting**: 登入 API 需要 rate limiting 防止暴力攻擊
4. **輸入驗證**: 所有輸入參數都需要驗證
5. **錯誤處理**: 不洩漏敏感資訊的錯誤訊息

### 隱私保護
- 僅儲存必要的使用者資訊
- 密碼使用 bcrypt 加密
- 第三方 token 不儲存在資料庫中

## 測試建議

### 單元測試
- Token 驗證邏輯
- 使用者建立/更新邏輯
- JWT 產生/驗證

### 整合測試
- 完整登入流程
- 帳號綁定邏輯
- API 端點測試

### 手動測試檢查清單
- [ ] Android Google 登入
- [ ] iOS Apple ID 登入
- [ ] iOS Google 登入
- [ ] Web Google 登入
- [ ] Web Apple ID 登入
- [ ] 相同 email 帳號自動綁定
- [ ] JWT token 有效性
- [ ] 登入後 API 呼叫權限

## 部署注意事項

### 環境變數設定
```env
# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Apple OAuth
APPLE_CLIENT_ID=your_apple_client_id
APPLE_TEAM_ID=your_apple_team_id
APPLE_KEY_ID=your_apple_key_id
APPLE_PRIVATE_KEY=your_apple_private_key

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_HOURS=168  # 7 days
```

### 第三方服務設定
1. **Google Cloud Console**: 設定 OAuth 2.0 用戶端 ID
2. **Apple Developer**: 設定 Sign in with Apple 服務
3. **Domain 驗證**: 確保所有 redirect URI 都已註冊

---

*此文件最後更新：2025年6月*