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

type CreateProductStockController struct {
	ProductID   string          `json:"product_id"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	StockAmount int             `json:"stock_amount"`
	Remaining   int             `json:"remaining"`
}

func (c *Controller) CreateProductStockController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`product_stocks.ctl.create.start`)

	var req CreateProductStockController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_stocks.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateProductStockService(ctx.Request.Context(), &CreateProductStockService{
		ProductID:   productID,
		UnitPrice:   req.UnitPrice,
		StockAmount: req.StockAmount,
		Remaining:   req.Remaining,
		MemberID:    memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`product_stocks.ctl.create.success`)
	base.Success(ctx, nil)
}
