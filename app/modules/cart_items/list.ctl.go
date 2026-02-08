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

type ListCartItemControllerRequest struct {
	base.RequestPaginate
}

type ListCartItemControllerResponses struct {
	ID              uuid.UUID       `json:"id"`
	CartID          uuid.UUID       `json:"cart_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (c *Controller) CartItemsList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListCartItemControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`cart_items.ctl.list.request`)

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

	data, page, err := c.svc.ListService(ctx, &ListCartItemServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`cart_items.ctl.list.callsvc`)

	var resp []*ListCartItemControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
