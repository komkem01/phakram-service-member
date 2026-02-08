package member_banks

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteMemberBankControllerRequest struct {
	ID string `uri:"id"`
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req DeleteMemberBankControllerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_banks.ctl.delete.request`)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if err := c.svc.DeleteService(ctx, id, memberID); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_banks.ctl.delete.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) MemberBanksDelete(ctx *gin.Context) {
	c.DeleteController(ctx)
}
