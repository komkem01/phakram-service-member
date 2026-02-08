package product_stocks

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

type UpdateProductStockControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateProductStockController struct {
	ProductID   string           `json:"product_id"`
	UnitPrice   *decimal.Decimal `json:"unit_price"`
	StockAmount *int             `json:"stock_amount"`
	Remaining   *int             `json:"remaining"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateProductStockControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_stocks.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateProductStockController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_stocks.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var productID uuid.UUID
	if req.ProductID != "" {
		productID, err = uuid.Parse(req.ProductID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateProductStockService{
		ProductID:   productID,
		UnitPrice:   req.UnitPrice,
		StockAmount: req.StockAmount,
		Remaining:   req.Remaining,
		MemberID:    memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`product_stocks.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) ProductStocksUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
