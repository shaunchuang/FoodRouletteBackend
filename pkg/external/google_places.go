package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// GooglePlacesService Google Places API 服務
type GooglePlacesService struct {
	apiKey string
	client *http.Client
}

// NewGooglePlacesService 建立 Google Places 服務
func NewGooglePlacesService(apiKey string) *GooglePlacesService {
	return &GooglePlacesService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// PlacesResponse Google Places API 回應結構
type PlacesResponse struct {
	Results []PlaceResult `json:"results"`
	Status  string        `json:"status"`
}

// PlaceResult 地點結果結構
type PlaceResult struct {
	PlaceID          string       `json:"place_id"`
	Name             string       `json:"name"`
	FormattedAddress string       `json:"formatted_address"`
	Geometry         Geometry     `json:"geometry"`
	Types            []string     `json:"types"`
	Rating           float64      `json:"rating"`
	PriceLevel       int          `json:"price_level"`
	Photos           []Photo      `json:"photos"`
	OpeningHours     OpeningHours `json:"opening_hours"`
}

// Geometry 地理位置結構
type Geometry struct {
	Location Location `json:"location"`
}

// Location 座標結構
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Photo 照片結構
type Photo struct {
	PhotoReference string `json:"photo_reference"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}

// OpeningHours 營業時間結構
type OpeningHours struct {
	OpenNow bool `json:"open_now"`
}

// PlaceDetailResponse 地點詳細資訊回應結構
type PlaceDetailResponse struct {
	Result PlaceDetail `json:"result"`
	Status string      `json:"status"`
}

// PlaceDetail 地點詳細資訊結構
type PlaceDetail struct {
	PlaceID              string       `json:"place_id"`
	Name                 string       `json:"name"`
	FormattedAddress     string       `json:"formatted_address"`
	FormattedPhoneNumber string       `json:"formatted_phone_number"`
	Geometry             Geometry     `json:"geometry"`
	Types                []string     `json:"types"`
	Rating               float64      `json:"rating"`
	PriceLevel           int          `json:"price_level"`
	Photos               []Photo      `json:"photos"`
	OpeningHours         OpeningHours `json:"opening_hours"`
	Website              string       `json:"website"`
}

// SearchNearbyRestaurants 搜尋附近餐廳
func (s *GooglePlacesService) SearchNearbyRestaurants(ctx context.Context, lat, lng float64, radius int) ([]domain.Restaurant, error) {
	baseURL := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"

	params := url.Values{}
	params.Set("location", fmt.Sprintf("%f,%f", lat, lng))
	params.Set("radius", fmt.Sprintf("%d", radius))
	params.Set("type", "restaurant")
	params.Set("key", s.apiKey)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		logger.Error("建立請求失敗", zap.Error(err))
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("Google Places API 請求失敗", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var placesResp PlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&placesResp); err != nil {
		logger.Error("解析 Google Places API 回應失敗", zap.Error(err))
		return nil, err
	}

	if placesResp.Status != "OK" {
		logger.Error("Google Places API 錯誤", zap.String("status", placesResp.Status))
		return nil, fmt.Errorf("Google Places API 錯誤: %s", placesResp.Status)
	}

	restaurants := make([]domain.Restaurant, 0, len(placesResp.Results))
	for _, place := range placesResp.Results {
		restaurant := s.convertToRestaurant(place)
		restaurants = append(restaurants, restaurant)
	}

	logger.Info("Google Places 搜尋成功",
		zap.Int("count", len(restaurants)),
		zap.Float64("lat", lat),
		zap.Float64("lng", lng),
		zap.Int("radius", radius),
	)

	return restaurants, nil
}

// GetRestaurantDetails 取得餐廳詳細資訊
func (s *GooglePlacesService) GetRestaurantDetails(ctx context.Context, googleID string) (*domain.Restaurant, error) {
	baseURL := "https://maps.googleapis.com/maps/api/place/details/json"

	params := url.Values{}
	params.Set("place_id", googleID)
	params.Set("fields", "place_id,name,formatted_address,formatted_phone_number,geometry,types,rating,price_level,photos,opening_hours,website")
	params.Set("key", s.apiKey)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		logger.Error("建立詳細資訊請求失敗", zap.Error(err))
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("Google Places 詳細資訊請求失敗", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var detailResp PlaceDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
		logger.Error("解析 Google Places 詳細資訊回應失敗", zap.Error(err))
		return nil, err
	}

	if detailResp.Status != "OK" {
		logger.Error("Google Places 詳細資訊 API 錯誤", zap.String("status", detailResp.Status))
		return nil, fmt.Errorf("Google Places 詳細資訊 API 錯誤: %s", detailResp.Status)
	}

	restaurant := s.convertDetailToRestaurant(detailResp.Result)

	logger.Info("Google Places 詳細資訊取得成功", zap.String("google_id", googleID))

	return &restaurant, nil
}

// convertToRestaurant 將 Google Places 結果轉換為餐廳實體
func (s *GooglePlacesService) convertToRestaurant(place PlaceResult) domain.Restaurant {
	var imageURL string
	if len(place.Photos) > 0 {
		imageURL = s.getPhotoURL(place.Photos[0].PhotoReference, 400)
	}

	// 判斷料理類型
	cuisine := s.determineCuisineType(place.Types)

	return domain.Restaurant{
		Name:       place.Name,
		Address:    place.FormattedAddress,
		Latitude:   place.Geometry.Location.Lat,
		Longitude:  place.Geometry.Location.Lng,
		Rating:     float32(place.Rating),
		PriceLevel: place.PriceLevel,
		Cuisine:    cuisine,
		GoogleID:   place.PlaceID,
		ImageURL:   imageURL,
		IsActive:   true,
	}
}

// convertDetailToRestaurant 將 Google Places 詳細資訊轉換為餐廳實體
func (s *GooglePlacesService) convertDetailToRestaurant(detail PlaceDetail) domain.Restaurant {
	var imageURL string
	if len(detail.Photos) > 0 {
		imageURL = s.getPhotoURL(detail.Photos[0].PhotoReference, 400)
	}

	// 判斷料理類型
	cuisine := s.determineCuisineType(detail.Types)

	return domain.Restaurant{
		Name:       detail.Name,
		Address:    detail.FormattedAddress,
		Phone:      detail.FormattedPhoneNumber,
		Latitude:   detail.Geometry.Location.Lat,
		Longitude:  detail.Geometry.Location.Lng,
		Rating:     float32(detail.Rating),
		PriceLevel: detail.PriceLevel,
		Cuisine:    cuisine,
		GoogleID:   detail.PlaceID,
		ImageURL:   imageURL,
		IsActive:   true,
	}
}

// getPhotoURL 取得照片 URL
func (s *GooglePlacesService) getPhotoURL(photoReference string, maxWidth int) string {
	return fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=%d&photoreference=%s&key=%s",
		maxWidth, photoReference, s.apiKey)
}

// determineCuisineType 根據 Google Places 類型判斷料理類型
func (s *GooglePlacesService) determineCuisineType(types []string) string {
	cuisineMap := map[string]string{
		"chinese_restaurant":    "中式料理",
		"japanese_restaurant":   "日式料理",
		"korean_restaurant":     "韓式料理",
		"italian_restaurant":    "義式料理",
		"american_restaurant":   "美式料理",
		"thai_restaurant":       "泰式料理",
		"indian_restaurant":     "印度料理",
		"mexican_restaurant":    "墨西哥料理",
		"french_restaurant":     "法式料理",
		"vietnamese_restaurant": "越南料理",
		"fast_food_restaurant":  "快餐",
		"pizza_restaurant":      "披薩",
		"seafood_restaurant":    "海鮮",
		"steakhouse":            "牛排",
		"bakery":                "烘焙",
		"cafe":                  "咖啡廳",
		"bar":                   "酒吧",
	}

	for _, placeType := range types {
		if cuisine, exists := cuisineMap[placeType]; exists {
			return cuisine
		}
	}

	return "餐廳"
}
