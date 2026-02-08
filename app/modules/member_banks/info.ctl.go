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

type InfoMemberBankControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoMemberBankControllerResponses struct {
	ID          uuid.UUID `json:"id"`
	MemberID    uuid.UUID `json:"member_id"`
	BankID      uuid.UUID `json:"bank_id"`
	BankNo      string    `json:"bank_no"`
	FirstnameTh string    `json:"firstname_th"`
	LastnameTh  string    `json:"lastname_th"`
	FirstnameEn string    `json:"firstname_en"`
	LastnameEn  string    `json:"lastname_en"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoMemberBankControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_banks.ctl.info.request`)

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
	span.AddEvent(`member_banks.ctl.info.callsvc`)

	var resp InfoMemberBankControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) MemberBanksInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
