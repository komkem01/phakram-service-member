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

type CreateMemberBankController struct {
	MemberID    string `json:"member_id"`
	BankID      string `json:"bank_id"`
	BankNo      string `json:"bank_no"`
	FirstnameTh string `json:"firstname_th"`
	LastnameTh  string `json:"lastname_th"`
	FirstnameEn string `json:"firstname_en"`
	LastnameEn  string `json:"lastname_en"`
	IsSystem    bool   `json:"is_system"`
	IsActive    *bool  `json:"is_active"`
}

func (c *Controller) CreateMemberBankController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`member_banks.ctl.create.start`)

	var req CreateMemberBankController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_banks.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}
	isAdmin := authmod.GetIsAdmin(ctx)
	if req.IsSystem && !isAdmin {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}
	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateMemberBankService(ctx.Request.Context(), &CreateMemberBankService{
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

	span.AddEvent(`member_banks.ctl.create.success`)
	base.Success(ctx, nil)
}
