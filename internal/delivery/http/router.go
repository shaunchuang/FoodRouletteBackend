package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaunchuang/food-roulette-backend/internal/delivery/http/handler"
	"github.com/shaunchuang/food-roulette-backend/internal/delivery/http/middleware"
	"github.com/shaunchuang/food-roulette-backend/internal/usecase"
)

// Router HTTP 路由器
type Router struct {
	userHandler       *handler.UserHandler
	restaurantHandler *handler.RestaurantHandler
	gameHandler       *handler.GameHandler
	adHandler         *handler.AdvertisementHandler
}

// NewRouter 建立新的路由器
func NewRouter(
	userHandler *handler.UserHandler,
	restaurantHandler *handler.RestaurantHandler,
	gameHandler *handler.GameHandler,
	adHandler *handler.AdvertisementHandler,
) *Router {
	return &Router{
		userHandler:       userHandler,
		restaurantHandler: restaurantHandler,
		gameHandler:       gameHandler,
		adHandler:         adHandler,
	}
}

// SetupRoutes 設定路由
func (r *Router) SetupRoutes(engine *gin.Engine, authService middleware.AuthService, userService *usecase.UserUseCase) {
	// 全域中介軟體
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.ErrorMiddleware())

	// 健康檢查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "美食沙漠樂園後端服務運行中",
		})
	})

	// API 版本群組
	v1 := engine.Group("/api/v1")
	{
		// 公開路由（不需要認證）
		public := v1.Group("/")
		{
			// 使用者認證相關
			auth := public.Group("/auth")
			{
				auth.POST("/register", r.userHandler.Register)
				auth.POST("/login", r.userHandler.Login)
			}

			// 餐廳相關（公開）
			restaurants := public.Group("/restaurants")
			{
				restaurants.GET("/search", r.restaurantHandler.SearchNearby)
				restaurants.GET("/:id", r.restaurantHandler.GetRestaurant)
			}

			// 廣告相關（公開瀏覽）
			ads := public.Group("/advertisements")
			{
				ads.GET("/", r.adHandler.GetActiveAds)
				ads.GET("/:id", r.adHandler.GetAdByID)
				ads.POST("/view", r.adHandler.RecordView)   // 記錄廣告瀏覽
				ads.POST("/click", r.adHandler.RecordClick) // 記錄廣告點擊
			}
		}

		// 需要認證的路由
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// 使用者相關
			users := protected.Group("/users")
			{
				users.GET("/profile", r.userHandler.GetProfile)
				users.PUT("/location", r.userHandler.UpdateLocation)
				users.GET("/location", r.userHandler.GetLocation)
			}

			// 最愛餐廳
			favorites := protected.Group("/favorites")
			{
				favorites.GET("/", r.restaurantHandler.GetFavorites)
				favorites.POST("/", r.restaurantHandler.AddToFavorites)
				favorites.DELETE("/:restaurant_id", r.restaurantHandler.RemoveFromFavorites)
			}

			// 遊戲相關
			games := protected.Group("/games")
			{
				games.POST("/start", r.gameHandler.StartGame)
				games.POST("/complete", r.gameHandler.CompleteGame)
				games.GET("/history", r.gameHandler.GetGameHistory)
			}

			// 廣告統計（需要認證）
			adStats := protected.Group("/advertisements")
			{
				adStats.GET("/:id/statistics", r.adHandler.GetStatistics)
			}
		}

		// 管理員路由（需要管理員權限）
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware(userService))
		{
			// 餐廳管理
			adminRestaurants := admin.Group("/restaurants")
			{
				adminRestaurants.POST("/", r.restaurantHandler.CreateRestaurant)
				adminRestaurants.GET("/", r.restaurantHandler.GetAllRestaurants)
				adminRestaurants.PUT("/:id", r.restaurantHandler.UpdateRestaurant)
			}

			// 廣告管理
			adminAds := admin.Group("/advertisements")
			{
				adminAds.POST("/", r.adHandler.CreateAd)
				adminAds.GET("/", r.adHandler.GetAllAds)
				adminAds.PUT("/:id", r.adHandler.UpdateAd)
				adminAds.DELETE("/:id", r.adHandler.DeleteAd)
			}
		}
	}
}
