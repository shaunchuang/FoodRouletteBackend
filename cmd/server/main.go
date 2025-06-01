package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL 驅動程式
	"github.com/shaunchuang/food-roulette-backend/internal/config"
	"github.com/shaunchuang/food-roulette-backend/internal/delivery/http"
	"github.com/shaunchuang/food-roulette-backend/internal/delivery/http/handler"
	"github.com/shaunchuang/food-roulette-backend/internal/repository/postgresql"
	"github.com/shaunchuang/food-roulette-backend/internal/usecase"
	"github.com/shaunchuang/food-roulette-backend/pkg/auth"
	"github.com/shaunchuang/food-roulette-backend/pkg/external"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("載入配置失敗:", err)
	}

	// 初始化日誌
	if err := logger.Init(cfg.Logger.Level); err != nil {
		log.Fatal("初始化日誌失敗:", err)
	}
	defer logger.Sync()

	logger.Info("美食沙漠樂園後端服務啟動中...")

	// 設定 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 建立 Gin 引擎
	engine := gin.New()

	// 初始化資料庫連接
	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("資料庫連接失敗", zap.Error(err))
	}
	defer db.Close()

	// 測試資料庫連接
	if err := db.Ping(); err != nil {
		logger.Fatal("資料庫連接測試失敗", zap.Error(err))
	}
	logger.Info("資料庫連接成功")

	// 初始化 Repositories
	userRepo := postgresql.NewUserRepository(db)
	restaurantRepo := postgresql.NewRestaurantRepository(db)
	favoriteRepo := postgresql.NewFavoriteRepository(db)
	gameRepo := postgresql.NewGameRepository(db)
	adRepo := postgresql.NewAdvertisementRepository(db)

	// 初始化 Services
	authService := auth.NewJWTService(cfg.Auth.Secret)
	var externalAPIService usecase.ExternalAPIService
	if cfg.GoogleAPI.PlacesAPIKey != "" {
		externalAPIService = external.NewGooglePlacesService(cfg.GoogleAPI.PlacesAPIKey)
	}

	// 初始化 Use Cases
	userUseCase := usecase.NewUserUseCase(userRepo, authService)
	restaurantUseCase := usecase.NewRestaurantUseCase(restaurantRepo, favoriteRepo, externalAPIService)
	gameUseCase := usecase.NewGameUseCase(gameRepo, restaurantRepo, favoriteRepo, adRepo)
	adUseCase := usecase.NewAdvertisementUseCase(adRepo)

	// 初始化 Handlers
	userHandler := handler.NewUserHandler(userUseCase)
	restaurantHandler := handler.NewRestaurantHandler(restaurantUseCase)
	gameHandler := handler.NewGameHandler(gameUseCase)
	adHandler := handler.NewAdvertisementHandler(adUseCase)

	// 初始化路由器
	router := http.NewRouter(userHandler, restaurantHandler, gameHandler, adHandler)
	router.SetupRoutes(engine, authService, userUseCase)

	// 啟動伺服器
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Info("伺服器啟動",
		zap.String("address", serverAddr),
		zap.String("mode", cfg.Server.Mode),
	)

	if err := engine.Run(serverAddr); err != nil {
		logger.Fatal("伺服器啟動失敗", zap.Error(err))
	}
}
