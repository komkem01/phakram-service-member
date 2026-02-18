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
	"member address not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบที่อยู่", nil, params...)
	},
	"member bank not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบบัญชีธนาคาร", nil, params...)
	},
	"member payment not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบข้อมูลการชำระเงิน", nil, params...)
	},
	"payment not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบข้อมูลการชำระเงิน", nil, params...)
	},
	"order is not pending": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "คำสั่งซื้อนี้ไม่อยู่ในสถานะรอชำระเงิน", nil, params...)
	},
	"payment confirmation already submitted": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ยืนยันการชำระเงินแล้ว อยู่ระหว่างรอตรวจสอบ", nil, params...)
	},
	"payment confirmation not submitted": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ยังไม่มีการยืนยันชำระเงินจากลูกค้า", nil, params...)
	},
	"rejection reason is required": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "กรุณาระบุเหตุผลที่ไม่อนุมัติ", nil, params...)
	},
	"shipping tracking number is required": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "กรุณากรอกเลขพัสดุก่อนเปลี่ยนเป็นสถานะกำลังจัดส่ง", nil, params...)
	},
	"payment was rejected waiting for resubmission": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่สามารถเปลี่ยนสถานะได้ กรุณารอลูกค้ายืนยันการชำระเงินใหม่", nil, params...)
	},
	"cannot cancel order after payment submission": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ยืนยันการชำระเงินแล้ว ไม่สามารถยกเลิกตรงได้ กรุณาติดต่อแอดมินเพื่อขอคืนเงิน", nil, params...)
	},
	"refund request requires payment submission": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "การขอคืนเงินทำได้หลังยืนยันการชำระเงินแล้วเท่านั้น", nil, params...)
	},
	"refund reason is required": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "กรุณาระบุเหตุผลที่ต้องการขอคืนเงิน", nil, params...)
	},
	"refund rejection reason is required": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "กรุณาระบุเหตุผลที่ปฏิเสธคำขอคืนเงิน", nil, params...)
	},
	"payment appeal reason is required": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "กรุณาระบุเหตุผลการอุทธรณ์การชำระเงิน", nil, params...)
	},
	"payment appeal is allowed only after rejection": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "อุทธรณ์ได้เฉพาะรายการที่ถูกปฏิเสธการชำระเงินแล้ว", nil, params...)
	},
	"payment is in use": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่สามารถลบได้ เนื่องจาก payment ถูกอ้างอิงอยู่", nil, params...)
	},
	"cart items not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบสินค้าในตะกร้า", nil, params...)
	},
	"cart item not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบสินค้าในตะกร้า", nil, params...)
	},
	"order not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบออเดอร์", nil, params...)
	},
	"order item not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบรายการสินค้าในออเดอร์", nil, params...)
	},
	"notification not found": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ไม่พบการแจ้งเตือน", nil, params...)
	},
	"order has no items": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "ออเดอร์นี้ไม่มีรายการสินค้าให้สั่งซ้ำ", nil, params...)
	},
	"reorder is allowed only for completed orders": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "สั่งซื้อซ้ำได้เฉพาะคำสั่งซื้อที่สำเร็จแล้ว", nil, params...)
	},
	"insufficient product stock": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "สต็อกสินค้าไม่เพียงพอสำหรับการสั่งซ้ำ", nil, params...)
	},
	"product is inactive": func(ctx *gin.Context, _ string, _ any, params ...map[string]string) error {
		return ValidateFailed(ctx, "มีสินค้าในออเดอร์ที่ปิดการขายแล้ว", nil, params...)
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
		return Unauthorized(ctx, "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง", nil, params...)
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
	"cart_items_cart_product_uidx":     "สินค้าในตะกร้าซ้ำ",
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
