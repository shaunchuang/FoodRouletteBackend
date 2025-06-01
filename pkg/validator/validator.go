package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// 註冊自定義驗證器
	validate.RegisterValidation("latitude", validateLatitude)
	validate.RegisterValidation("longitude", validateLongitude)
}

// ValidateStruct 驗證結構體
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// GetValidationErrors 取得詳細的驗證錯誤訊息
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := getJSONFieldName(fieldError)
			errors[fieldName] = getErrorMessage(fieldError)
		}
	}

	return errors
}

// getJSONFieldName 取得 JSON 欄位名稱
func getJSONFieldName(fieldError validator.FieldError) string {
	field := fieldError.Field()

	// 嘗試從 struct tag 中取得 json 名稱
	if fieldError.StructNamespace() != "" {
		parts := strings.Split(fieldError.StructNamespace(), ".")
		if len(parts) > 1 {
			field = strings.ToLower(parts[len(parts)-1])
		}
	}

	return field
}

// getErrorMessage 根據驗證規則取得錯誤訊息
func getErrorMessage(fieldError validator.FieldError) string {
	field := fieldError.Field()
	tag := fieldError.Tag()
	param := fieldError.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s 欄位為必填", field)
	case "email":
		return fmt.Sprintf("%s 必須是有效的電子郵件格式", field)
	case "min":
		return fmt.Sprintf("%s 最少需要 %s 個字元", field, param)
	case "max":
		return fmt.Sprintf("%s 最多只能 %s 個字元", field, param)
	case "latitude":
		return fmt.Sprintf("%s 必須是有效的緯度值 (-90 到 90)", field)
	case "longitude":
		return fmt.Sprintf("%s 必須是有效的經度值 (-180 到 180)", field)
	case "gte":
		return fmt.Sprintf("%s 必須大於或等於 %s", field, param)
	case "lte":
		return fmt.Sprintf("%s 必須小於或等於 %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s 必須是以下值之一: %s", field, param)
	default:
		return fmt.Sprintf("%s 欄位驗證失敗", field)
	}
}

// validateLatitude 自定義緯度驗證器
func validateLatitude(fl validator.FieldLevel) bool {
	lat := fl.Field().Float()
	return lat >= -90.0 && lat <= 90.0
}

// validateLongitude 自定義經度驗證器
func validateLongitude(fl validator.FieldLevel) bool {
	lng := fl.Field().Float()
	return lng >= -180.0 && lng <= 180.0
}

// ValidateEmail 驗證電子郵件格式
func ValidateEmail(email string) bool {
	return validate.Var(email, "required,email") == nil
}

// ValidatePassword 驗證密碼強度
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("密碼至少需要 6 個字元")
	}
	if len(password) > 128 {
		return fmt.Errorf("密碼最多只能 128 個字元")
	}
	return nil
}

// ValidateLocation 驗證地理位置
func ValidateLocation(lat, lng float64) error {
	if lat < -90.0 || lat > 90.0 {
		return fmt.Errorf("緯度必須在 -90 到 90 之間")
	}
	if lng < -180.0 || lng > 180.0 {
		return fmt.Errorf("經度必須在 -180 到 180 之間")
	}
	return nil
}
