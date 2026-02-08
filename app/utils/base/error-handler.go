package base

import (
	"github.com/gin-gonic/gin"
)

// ResponseFunction defines a function that responds to HTTP requests
type ResponseFunction func(ctx *gin.Context, message string, data any, params ...map[string]string) error

// Use string literal keys for error mappings (Go compiler will detect duplicate keys)
var errorMappings = map[string]ResponseFunction{
	"sql: no rows in result set": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบข้อมูล", nil, params...)
	},
	"no rows in result set": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบข้อมูล", nil, params...)
	},
	"email already exists": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "อีเมลซ้ำ", nil, params...)
	},
	"phone already exists": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "เบอร์โทรซ้ำ", nil, params...)
	},
	"cart not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบตะกร้า", nil, params...)
	},
	"address not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบที่อยู่", nil, params...)
	},
	"cart items not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบสินค้าในตะกร้า", nil, params...)
	},
	"product stock not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบสต็อกสินค้า", nil, params...)
	},
	"default tier not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบระดับทั่วไป", nil, params...)
	},
	"default status not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบสถานะใช้งาน", nil, params...)
	},
}

func HandleError(ctx *gin.Context, err error) {
	if err != nil {
		// Use map to lookup error response function (convert error to string for lookup)
		if responseFunc, ok := errorMappings[err.Error()]; ok {
			responseFunc(ctx, err.Error(), nil)
			return
		}

		// If error doesn't match any mapping, return the original error
		InternalServerError(ctx, err.Error(), nil)
	}
}
