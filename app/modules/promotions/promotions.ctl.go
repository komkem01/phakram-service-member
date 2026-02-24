package promotions

import (
	"fmt"
	"strconv"
	"strings"

	"phakram/app/modules/auth"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type ListPromotionsControllerRequest struct {
	base.RequestPaginate
	IsActive *bool `form:"is_active"`
}

type CreatePromotionControllerRequest struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	DiscountType   string   `json:"discount_type"`
	DiscountValue  float64  `json:"discount_value"`
	MaxDiscount    *float64 `json:"max_discount"`
	MinOrderAmount float64  `json:"min_order_amount"`
	UsageLimit     *int     `json:"usage_limit"`
	UsagePerMember *int     `json:"usage_per_member"`
	StartsAt       *string  `json:"starts_at"`
	EndsAt         *string  `json:"ends_at"`
	IsActive       *bool    `json:"is_active"`
}

type UpdatePromotionControllerRequest struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	DiscountType   string   `json:"discount_type"`
	DiscountValue  float64  `json:"discount_value"`
	MaxDiscount    *float64 `json:"max_discount"`
	MinOrderAmount float64  `json:"min_order_amount"`
	UsageLimit     *int     `json:"usage_limit"`
	UsagePerMember *int     `json:"usage_per_member"`
	StartsAt       *string  `json:"starts_at"`
	EndsAt         *string  `json:"ends_at"`
	IsActive       *bool    `json:"is_active"`
}

type ValidatePromotionControllerRequest struct {
	Code        string  `json:"code"`
	OrderAmount float64 `json:"order_amount"`
}

type UsePromotionControllerRequest struct {
	OrderID        *string `json:"order_id"`
	DiscountAmount any     `json:"discount_amount"`
}

func parseDiscountAmount(raw any) (float64, error) {
	switch value := raw.(type) {
	case float64:
		return value, nil
	case float32:
		return float64(value), nil
	case int:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("unsupported discount_amount type")
	}
}

type ListMemberPromotionsControllerRequest struct {
	base.RequestPaginate
}

type ListPromotionUsagesControllerRequest struct {
	base.RequestPaginate
	PromotionID string `form:"promotion_id"`
}

func (c *Controller) ListController(ctx *gin.Context) {
	var req ListPromotionsControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.List(ctx.Request.Context(), &ListPromotionsServiceRequest{
		RequestPaginate: req.RequestPaginate,
		IsActive:        req.IsActive,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Paginate(ctx, data, page)
}

func (c *Controller) InfoController(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		base.BadRequest(ctx, "ไม่พบรหัสโปรโมชั่น", nil)
		return
	}

	data, err := c.svc.Info(ctx.Request.Context(), id)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, data)
}

func (c *Controller) CreateController(ctx *gin.Context) {
	var req CreatePromotionControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		base.BadRequest(ctx, "กรุณาระบุข้อมูลให้ครบถ้วน", nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if err := c.svc.Create(ctx.Request.Context(), &CreatePromotionServiceRequest{
		Code:           req.Code,
		Name:           req.Name,
		Description:    req.Description,
		DiscountType:   req.DiscountType,
		DiscountValue:  req.DiscountValue,
		MaxDiscount:    req.MaxDiscount,
		MinOrderAmount: req.MinOrderAmount,
		UsageLimit:     req.UsageLimit,
		UsagePerMember: req.UsagePerMember,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		IsActive:       isActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		base.BadRequest(ctx, "ไม่พบรหัสโปรโมชั่น", nil)
		return
	}

	var req UpdatePromotionControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		base.BadRequest(ctx, "กรุณาระบุข้อมูลให้ครบถ้วน", nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if err := c.svc.Update(ctx.Request.Context(), &UpdatePromotionServiceRequest{
		ID:             id,
		Code:           req.Code,
		Name:           req.Name,
		Description:    req.Description,
		DiscountType:   req.DiscountType,
		DiscountValue:  req.DiscountValue,
		MaxDiscount:    req.MaxDiscount,
		MinOrderAmount: req.MinOrderAmount,
		UsageLimit:     req.UsageLimit,
		UsagePerMember: req.UsagePerMember,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		IsActive:       isActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		base.BadRequest(ctx, "ไม่พบรหัสโปรโมชั่น", nil)
		return
	}

	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) ValidateController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var req ValidatePromotionControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.Validate(ctx.Request.Context(), memberID, &ValidatePromotionServiceRequest{
		Code:        req.Code,
		OrderAmount: req.OrderAmount,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, data)
}

func (c *Controller) UseController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	promotionID := strings.TrimSpace(ctx.Param("id"))
	if promotionID == "" {
		base.BadRequest(ctx, "ไม่พบรหัสโปรโมชั่น", nil)
		return
	}

	var req UsePromotionControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	discountAmount, err := parseDiscountAmount(req.DiscountAmount)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.Use(ctx.Request.Context(), &UsePromotionServiceRequest{
		PromotionID:    promotionID,
		MemberID:       memberID,
		OrderID:        req.OrderID,
		DiscountAmount: discountAmount,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) ListAvailableForMemberController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var req ListMemberPromotionsControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListAvailableForMember(ctx.Request.Context(), memberID, &ListMemberPromotionsServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Paginate(ctx, data, page)
}

func (c *Controller) ListMyController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var req ListMemberPromotionsControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListMy(ctx.Request.Context(), memberID, &ListMemberPromotionsServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Paginate(ctx, data, page)
}

func (c *Controller) CollectController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	promotionID := strings.TrimSpace(ctx.Param("id"))
	if promotionID == "" {
		base.BadRequest(ctx, "ไม่พบรหัสโปรโมชั่น", nil)
		return
	}

	if err := c.svc.Collect(ctx.Request.Context(), memberID, promotionID); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) ReportSummaryController(ctx *gin.Context) {
	data, err := c.svc.ReportSummary(ctx.Request.Context())
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, data)
}

func (c *Controller) ListUsagesController(ctx *gin.Context) {
	var req ListPromotionUsagesControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListUsages(ctx.Request.Context(), &ListPromotionUsagesServiceRequest{
		RequestPaginate: req.RequestPaginate,
		PromotionID:     req.PromotionID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Paginate(ctx, data, page)
}
