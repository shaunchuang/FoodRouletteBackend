package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/config"
	"github.com/shaunchuang/food-roulette-backend/internal/delivery/http"
	"github.com/shaunchuang/food-roulette-backend/internal/delivery/http/handler"
	"github.com/shaunchuang/food-roulette-backend/internal/usecase"
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

	// TODO: 初始化資料庫連接
	// db, err := sql.Open("postgres", cfg.Database.GetDSN())
	// if err != nil {
	//     logger.Fatal("資料庫連接失敗", zap.Error(err))
	// }
	// defer db.Close()

	// TODO: 初始化 Repositories
	// userRepo := postgresql.NewUserRepository(db)
	// restaurantRepo := postgresql.NewRestaurantRepository(db)
	// favoriteRepo := postgresql.NewFavoriteRepository(db)
	// gameRepo := postgresql.NewGameRepository(db)
	// adRepo := postgresql.NewAdvertisementRepository(db)

	// TODO: 初始化 Services
	// authService := auth.NewJWTService(cfg.Auth.Secret)
	// externalAPIService := external.NewGooglePlacesService(cfg.GoogleAPI.Key)

	// TODO: 初始化 Use Cases
	// userUseCase := usecase.NewUserUseCase(userRepo, authService)
	// restaurantUseCase := usecase.NewRestaurantUseCase(restaurantRepo, favoriteRepo, externalAPIService)
	// gameUseCase := usecase.NewGameUseCase(gameRepo, restaurantRepo, favoriteRepo, adRepo)

	// 暫時建立空的 Use Cases（之後會實作真正的依賴注入）
	var userUseCase *usecase.UserUseCase
	var restaurantUseCase *usecase.RestaurantUseCase
	var gameUseCase *usecase.GameUseCase

	// 初始化 Handlers
	userHandler := handler.NewUserHandler(userUseCase)
	restaurantHandler := handler.NewRestaurantHandler(restaurantUseCase)
	gameHandler := handler.NewGameHandler(gameUseCase)
	adHandler := handler.NewAdvertisementHandler()

	// 初始化路由器
	router := http.NewRouter(userHandler, restaurantHandler, gameHandler, adHandler)
	router.SetupRoutes(engine)

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