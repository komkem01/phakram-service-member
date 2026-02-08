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

type UpdateOrderControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateOrderController struct {
	OrderNo        string           `json:"order_no"`
	MemberID       string           `json:"member_id"`
	PaymentID      string           `json:"payment_id"`
	AddressID      string           `json:"address_id"`
	Status         string           `json:"status"`
	TotalAmount    *decimal.Decimal `json:"total_amount"`
	DiscountAmount *decimal.Decimal `json:"discount_amount"`
	NetAmount      *decimal.Decimal `json:"net_amount"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateOrderControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`orders.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateOrderController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`orders.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
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

	if err := c.svc.UpdateService(ctx, id, &UpdateOrderService{
		OrderNo:        req.OrderNo,
		MemberID:       memberID,
		PaymentID:      paymentID,
		AddressID:      addressID,
		Status:         req.Status,
		TotalAmount:    req.TotalAmount,
		DiscountAmount: req.DiscountAmount,
		NetAmount:      req.NetAmount,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`orders.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) OrdersUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
