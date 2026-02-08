package order_items

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

type CreateOrderItemController struct {
	OrderID         string          `json:"order_id"`
	ProductID       string          `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
}

func (c *Controller) CreateOrderItemController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`order_items.ctl.create.start`)

	var req CreateOrderItemController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`order_items.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateOrderItemService(ctx.Request.Context(), &CreateOrderItemService{
		OrderID:         orderID,
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
		MemberID:        memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`order_items.ctl.create.success`)
	base.Success(ctx, nil)
}
