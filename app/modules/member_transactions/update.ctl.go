package member_transactions

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateMemberTransactionControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateMemberTransactionController struct {
	MemberID string `json:"member_id"`
	Action   string `json:"action"`
	Details  string `json:"details"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateMemberTransactionControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_transactions.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateMemberTransactionController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_transactions.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateMemberTransactionService{
		MemberID: memberID,
		Action:   req.Action,
		Details:  req.Details,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_transactions.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) MemberTransactionsUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
