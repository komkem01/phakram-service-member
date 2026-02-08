package orders

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateOrderController struct {
	MemberID       string          `json:"member_id"`
	CartID         string          `json:"cart_id"`
	CartItemIDs    []string        `json:"cart_item_ids"`
	OrderNo        string          `json:"order_no"`
	PaymentID      string          `json:"payment_id"`
	AddressID      string          `json:"address_id"`
	Address        *OrderAddress   `json:"address"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
}

type OrderAddress struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	IsDefault     bool   `json:"is_default"`
	AddressNo     string `json:"address_no"`
	Village       string `json:"village"`
	Alley         string `json:"alley"`
	SubDistrictID string `json:"sub_district_id"`
	DistrictID    string `json:"district_id"`
	ProvinceID    string `json:"province_id"`
	ZipcodeID     string `json:"zipcode_id"`
}

func (c *Controller) CreateOrderController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.create.start`)

	var req CreateOrderController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`orders.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if authmod.GetIsAdmin(ctx) && req.MemberID != "" {
		adminMemberID, err := uuid.Parse(req.MemberID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		memberID = adminMemberID
	}

	cartID, err := uuid.Parse(req.CartID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var paymentID uuid.UUID
	if req.PaymentID != "" {
		paymentID, err = uuid.Parse(req.PaymentID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	var addressID uuid.UUID
	if req.AddressID != "" {
		addressID, err = uuid.Parse(req.AddressID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	if addressID == uuid.Nil && req.Address == nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	itemIDs := make([]uuid.UUID, 0, len(req.CartItemIDs))
	for _, rawID := range req.CartItemIDs {
		parsed, err := uuid.Parse(rawID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		itemIDs = append(itemIDs, parsed)
	}

	var addressPayload *CreateOrderAddress
	if req.Address != nil {
		addressPayload = &CreateOrderAddress{
			FirstName:     req.Address.FirstName,
			LastName:      req.Address.LastName,
			Phone:         req.Address.Phone,
			IsDefault:     req.Address.IsDefault,
			AddressNo:     req.Address.AddressNo,
			Village:       req.Address.Village,
			Alley:         req.Address.Alley,
			SubDistrictID: req.Address.SubDistrictID,
			DistrictID:    req.Address.DistrictID,
			ProvinceID:    req.Address.ProvinceID,
			ZipcodeID:     req.Address.ZipcodeID,
		}
	}

	if err := c.svc.CreateOrderService(ctx.Request.Context(), &CreateOrderService{
		CartID:         cartID,
		CartItemIDs:    itemIDs,
		OrderNo:        req.OrderNo,
		MemberID:       memberID,
		PaymentID:      paymentID,
		AddressID:      addressID,
		Address:        addressPayload,
		DiscountAmount: req.DiscountAmount,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.create.success`)
	base.Success(ctx, nil)
}
