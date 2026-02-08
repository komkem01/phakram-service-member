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

type InfoMemberAccountControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoMemberAccountControllerRequestQuery struct {
	Email string `form:"email"`
}

type InfoMemberAccountControllerRequestByMemberID struct {
	MemberID string `uri:"member_id"`
}

type InfoMemberAccountControllerResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoMemberAccountControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_accounts.ctl.info.request`)

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
	span.AddEvent(`member_accounts.ctl.info.callsvc`)

	var resp InfoMemberAccountControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) MemberAccountsInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}

func (c *Controller) MemberAccountsInfoByEmail(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoMemberAccountControllerRequestQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	if req.Email == "" {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_accounts.ctl.info_by_email.request`)

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

	data, err := c.svc.InfoByEmailService(ctx, req.Email, memberID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_accounts.ctl.info_by_email.callsvc`)

	var resp InfoMemberAccountControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) MemberAccountsInfoByMemberID(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	var req InfoMemberAccountControllerRequestByMemberID
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_accounts.ctl.info_by_member_id.request`)

	memberIDParam, err := uuid.Parse(req.MemberID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	requesterMemberID := uuid.Nil
	isAdmin := authmod.GetIsAdmin(ctx)
	if !isAdmin {
		var ok bool
		requesterMemberID, ok = authmod.GetMemberID(ctx)
		if !ok {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			return
		}
	}

	data, err := c.svc.InfoByMemberIDService(ctx, memberIDParam, requesterMemberID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_accounts.ctl.info_by_member_id.callsvc`)

	var resp InfoMemberAccountControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}
	base.Success(ctx, resp)
}
