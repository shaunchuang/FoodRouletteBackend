package domain

import "errors"

// 通用錯誤
var (
	ErrNotFound      = errors.New("資源不存在")
	ErrUnauthorized  = errors.New("未授權")
	ErrForbidden     = errors.New("權限不足")
	ErrInvalidInput  = errors.New("輸入資料無效")
	ErrConflict      = errors.New("資源衝突")
	ErrInternalError = errors.New("內部錯誤")
)

// 使用者相關錯誤
var (
	ErrUserNotFound    = errors.New("使用者不存在")
	ErrUserExists      = errors.New("使用者已存在")
	ErrInvalidPassword = errors.New("密碼錯誤")
	ErrInvalidEmail    = errors.New("電子郵件格式錯誤")
	ErrWeakPassword    = errors.New("密碼強度不足")
)

// 餐廳相關錯誤
var (
	ErrRestaurantNotFound = errors.New("餐廳不存在")
	ErrInvalidLocation    = errors.New("地理位置無效")
	ErrInvalidRadius      = errors.New("搜尋半徑無效")
)

// 最愛餐廳相關錯誤
var (
	ErrFavoriteExists   = errors.New("餐廳已在最愛清單中")
	ErrFavoriteNotFound = errors.New("最愛餐廳不存在")
)

// 遊戲相關錯誤
var (
	ErrGameSessionNotFound = errors.New("遊戲會話不存在")
	ErrGameSessionExpired  = errors.New("遊戲會話已過期")
	ErrInvalidGameType     = errors.New("無效的遊戲類型")
	ErrGameAlreadyComplete = errors.New("遊戲已完成")
)

// 廣告相關錯誤
var (
	ErrAdvertisementNotFound = errors.New("廣告不存在")
	ErrAdvertisementExpired  = errors.New("廣告已過期")
	ErrInvalidPeriod         = errors.New("無效的時間週期")
	ErrAdViewCooldown        = errors.New("廣告瀏覽冷卻時間未到")
	ErrAdClickCooldown       = errors.New("廣告點擊冷卻時間未到")
)

// JWT 相關錯誤
var (
	ErrInvalidToken = errors.New("無效的 token")
	ErrExpiredToken = errors.New("token 已過期")
	ErrMissingToken = errors.New("缺少 token")
)

// 外部 API 錯誤
var (
	ErrExternalAPIFailed = errors.New("外部 API 請求失敗")
	ErrGoogleAPIFailed   = errors.New("Google API 請求失敗")
	ErrAPIQuotaExceeded  = errors.New("API 配額已用完")
)

// 檔案上傳錯誤
var (
	ErrFileTooBig      = errors.New("檔案過大")
	ErrInvalidFileType = errors.New("不支援的檔案類型")
	ErrUploadFailed    = errors.New("檔案上傳失敗")
)

// 速率限制錯誤
var (
	ErrRateLimitExceeded = errors.New("請求頻率過高")
	ErrTooManyRequests   = errors.New("請求次數過多")
)
