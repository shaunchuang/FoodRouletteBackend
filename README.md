# 美食沙漠樂園後端 (Food Roulette Backend)

這是「美食沙漠樂園」應用程式的後端服務，採用 Go 語言和 Gin 框架開發，遵循 Clean Architecture 架構模式。

## 專案架構

本專案採用 **Clean Architecture + Domain-Driven Design (DDD)** 的混合架構：

```
FoodRouletteBackend/
├── cmd/
│   └── server/                 # 應用程式入口點
├── internal/
│   ├── config/                 # 配置管理
│   ├── domain/                 # 領域實體 (DDD)
│   ├── usecase/                # 業務邏輯層 (Clean Architecture)
│   ├── delivery/
│   │   └── http/              # HTTP 交付層
│   │       ├── handler/       # HTTP 處理器
│   │       └── middleware/    # 中介軟體
│   └── repository/
│       └── postgresql/        # 資料持久層實作
├── pkg/                       # 公共套件
│   ├── logger/               # 日誌工具
│   └── validator/            # 驗證工具
├── api/                      # API 文件
├── migrations/               # 資料庫遷移檔案
└── .env.example             # 環境變數範例
```

## 核心功能

### 🎮 遊戲化推薦
- **餐廳輪盤**：隨機選擇餐廳的輪盤遊戲
- **多種玩法**：支援骰子、塔羅牌、拼圖等多種遊戲模式
- **廣告整合**：遊戲過程中展示餐廳廣告

### 👤 使用者管理
- 使用者註冊與登入
- 個人資料管理
- 位置資訊更新

### 🍽️ 餐廳功能
- 根據 GPS 定位搜尋附近餐廳
- 個人最愛餐廳清單
- 餐廳詳細資訊查看

### 📊 廣告系統
- 廣告展示與點擊追蹤
- 廣告效果統計
- 廣告管理後台

## 技術棧

- **語言**: Go 1.21+
- **框架**: Gin (HTTP Web Framework)
- **資料庫**: PostgreSQL
- **遷移工具**: golang-migrate
- **日誌**: Zap (Uber 的高效能日誌庫)
- **架構**: Clean Architecture + DDD

## 快速開始

### 1. 環境需求

- Go 1.21 或更高版本
- PostgreSQL 12 或更高版本

### 2. 安裝依賴

```bash
go mod download
```

### 3. 設定環境變數

複製環境變數範例檔案：

```bash
cp .env.example .env
```

編輯 `.env` 檔案，設定您的資料庫連接資訊：

```env
# 資料庫配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=food_roulette
```

### 4. 建立資料庫

```bash
createdb food_roulette
```

### 5. 執行資料庫遷移

```bash
migrate -path migrations -database "postgres://username:password@localhost/food_roulette?sslmode=disable" up
```

### 6. 啟動伺服器

```bash
go run cmd/server/main.go
```

伺服器將在 `http://localhost:8080` 啟動。

## API 端點

### 健康檢查
- `GET /health` - 伺服器健康狀態

### 認證
- `POST /api/v1/auth/register` - 使用者註冊
- `POST /api/v1/auth/login` - 使用者登入

### 使用者
- `GET /api/v1/users/profile` - 取得使用者資料
- `PUT /api/v1/users/location` - 更新使用者位置

### 餐廳
- `GET /api/v1/restaurants/search` - 搜尋附近餐廳
- `GET /api/v1/restaurants/:id` - 取得餐廳詳細資訊

### 最愛餐廳
- `GET /api/v1/favorites` - 取得最愛餐廳清單
- `POST /api/v1/favorites` - 新增最愛餐廳
- `DELETE /api/v1/favorites/:restaurant_id` - 移除最愛餐廳

### 遊戲
- `POST /api/v1/games/start` - 開始遊戲
- `POST /api/v1/games/complete` - 完成遊戲
- `GET /api/v1/games/history` - 取得遊戲歷史

### 廣告
- `GET /api/v1/advertisements` - 取得活躍廣告
- `GET /api/v1/advertisements/:id/statistics` - 取得廣告統計

## 開發指南

### 專案結構說明

1. **Domain Layer** (`internal/domain/`): 包含業務實體和值物件
2. **Use Case Layer** (`internal/usecase/`): 包含應用程式業務邏輯
3. **Interface Adapters** (`internal/delivery/`, `internal/repository/`): 處理外部介面
4. **Frameworks & Drivers** (`cmd/`, `pkg/`): 框架和驅動程式

### 新增功能

1. 在 `internal/domain/` 中定義新的實體
2. 在 `internal/usecase/interfaces.go` 中定義新的介面
3. 在 `internal/usecase/` 中實作業務邏輯
4. 在 `internal/delivery/http/handler/` 中實作 HTTP 處理器
5. 在 `internal/delivery/http/router.go` 中註冊新的路由

## 資料庫設計

資料庫包含以下主要資料表：

- `users` - 使用者資訊
- `user_locations` - 使用者位置
- `restaurants` - 餐廳資訊
- `favorite_restaurants` - 最愛餐廳
- `game_sessions` - 遊戲會話
- `advertisements` - 廣告資訊
- `ad_views` / `ad_clicks` - 廣告統計

## 貢獻指南

1. Fork 這個專案
2. 建立您的功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的變更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 開啟一個 Pull Request

## 授權

本專案採用 MIT 授權 - 詳見 [LICENSE](LICENSE) 檔案。

## 聯絡資訊

如有任何問題或建議，請透過 Issue 與我們聯絡。