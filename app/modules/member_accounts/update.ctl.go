package member_accounts

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateMemberAccountControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateMemberAccountController struct {
	MemberID string `json:"member_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateMemberAccountControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_accounts.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateMemberAccountController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_accounts.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateMemberAccountService{
		MemberID: memberID,
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_accounts.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) MemberAccountsUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
