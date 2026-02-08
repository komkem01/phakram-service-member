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

type UpdateMemberBankControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateMemberBankController struct {
	MemberID    string `json:"member_id"`
	BankID      string `json:"bank_id"`
	BankNo      string `json:"bank_no"`
	FirstnameTh string `json:"firstname_th"`
	LastnameTh  string `json:"lastname_th"`
	FirstnameEn string `json:"firstname_en"`
	LastnameEn  string `json:"lastname_en"`
	IsSystem    *bool  `json:"is_system"`
	IsActive    *bool  `json:"is_active"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateMemberBankControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_banks.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateMemberBankController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_banks.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}
	isAdmin := authmod.GetIsAdmin(ctx)
	if req.IsSystem != nil && !isAdmin {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}
	var bankID uuid.UUID
	if req.BankID != "" {
		bankID, err = uuid.Parse(req.BankID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateMemberBankService{
		MemberID:    memberID,
		BankID:      bankID,
		BankNo:      req.BankNo,
		FirstnameTh: req.FirstnameTh,
		LastnameTh:  req.LastnameTh,
		FirstnameEn: req.FirstnameEn,
		LastnameEn:  req.LastnameEn,
		IsSystem:    req.IsSystem,
		IsActive:    req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_banks.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) MemberBanksUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
