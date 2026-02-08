package product_stocks

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoProductStockControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoProductStockControllerResponses struct {
	ID          uuid.UUID       `json:"id"`
	ProductID   uuid.UUID       `json:"product_id"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	StockAmount int             `json:"stock_amount"`
	Remaining   int             `json:"remaining"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoProductStockControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_stocks.ctl.info.request`)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.InfoService(ctx, id)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`product_stocks.ctl.info.callsvc`)

	var resp InfoProductStockControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) ProductStocksInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
