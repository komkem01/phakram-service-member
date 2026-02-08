package cart_items

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SelectCartItemsControllerRequest struct {
	CartID      string   `json:"cart_id"`
	CartItemIDs []string `json:"cart_item_ids"`
}

func (c *Controller) SelectItemsController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`cart_items.ctl.select.start`)

	var req SelectCartItemsControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	cartID, err := uuid.Parse(req.CartID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
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

	if err := c.svc.SelectItemsService(ctx.Request.Context(), &SelectCartItemsServiceRequest{
		CartID:   cartID,
		ItemIDs:  itemIDs,
		MemberID: memberID,
		IsAdmin:  isAdmin,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`cart_items.ctl.select.success`)
	base.Success(ctx, nil)
}
