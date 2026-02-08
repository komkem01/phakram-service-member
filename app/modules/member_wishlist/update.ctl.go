package member_wishlist

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

type UpdateMemberWishlistControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateMemberWishlistController struct {
	MemberID        string           `json:"member_id"`
	ProductID       string           `json:"product_id"`
	Quantity        *int             `json:"quantity"`
	PricePerUnit    *decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount *decimal.Decimal `json:"total_item_amount"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateMemberWishlistControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_wishlist.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateMemberWishlistController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_wishlist.ctl.update.request_body`)

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

	if err := c.svc.UpdateService(ctx, id, &UpdateMemberWishlistService{
		MemberID:        memberID,
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_wishlist.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) MemberWishlistUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
