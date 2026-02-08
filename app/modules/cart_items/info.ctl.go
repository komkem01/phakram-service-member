package cart_items

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

type InfoCartItemControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoCartItemControllerResponses struct {
	ID              uuid.UUID       `json:"id"`
	CartID          uuid.UUID       `json:"cart_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	IsSelected      bool            `json:"is_selected"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoCartItemControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`cart_items.ctl.info.request`)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	memberID := uuid.Nil
	isAdmin := authmod.GetIsAdmin(ctx)
	if !isAdmin {
		var ok bool
		memberID, ok = authmod.GetMemberID(ctx)
		if !ok {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			return
		}
	}

	data, err := c.svc.InfoService(ctx, id, memberID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`cart_items.ctl.info.callsvc`)

	var resp InfoCartItemControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) CartItemsInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
