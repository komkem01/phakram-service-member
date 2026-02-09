package base

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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
	"invalid credentials": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return Unauthorized(ctx, "ข้อมูลรับรองไม่ถูกต้อง", nil, params...)
	},
	"token expired": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return Unauthorized(ctx, "โทเค็นหมดอายุ", nil, params...)
	},
	"invalid token": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return Unauthorized(ctx, "โทเค็นไม่ถูกต้อง", nil, params...)
	},
	"forbidden": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return Forbidden(ctx, "ไม่มีสิทธิ์เข้าถึง", nil, params...)
	},

	"prefix-not-found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบคำนำหน้า", nil, params...)
	},
	"gender-not-found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบเพศ", nil, params...)
	},
}

var duplicateConstraintMessages = map[string]string{
	"provinces_name_uidx":              "ชื่อจังหวัดซ้ำ",
	"districts_province_name_uidx":     "ชื่ออำเภอซ้ำ",
	"sub_districts_district_name_uidx": "ชื่อตำบลซ้ำ",
	"zipcodes_sub_district_name_uidx":  "รหัสไปรษณีย์ซ้ำ",
	"members_member_no_uidx":           "รหัสสมาชิกซ้ำ",
	"members_phone_uidx":               "เบอร์โทรซ้ำ",
}

func duplicateErrorMessage(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return "", false
	}
	if pgErr.Code != "23505" {
		return "", false
	}
	if msg, ok := duplicateConstraintMessages[pgErr.ConstraintName]; ok {
		return msg, true
	}
	return "ข้อมูลซ้ำ", true
}

func HandleError(ctx *gin.Context, err error) {
	if err != nil {
		if msg, ok := duplicateErrorMessage(err); ok {
			ValidateFailed(ctx, msg, nil)
			return
		}

		// Use map to lookup error response function (convert error to string for lookup)
		if responseFunc, ok := errorMappings[err.Error()]; ok {
			responseFunc(ctx, err.Error(), nil)
			return
		}

		// If error doesn't match any mapping, return the original error
		InternalServerError(ctx, err.Error(), nil)
	}
}
