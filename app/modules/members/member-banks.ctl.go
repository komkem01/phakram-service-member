package members

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberBankURIRequest struct {
	MemberID string `uri:"id"`
	BankID   string `uri:"bank_id"`
}

type CreateMemberBankControllerRequest struct {
	BankID      string `json:"bank_id"`
	BankNo      string `json:"bank_no"`
	FirstnameTh string `json:"firstname_th"`
	LastnameTh  string `json:"lastname_th"`
	FirstnameEn string `json:"firstname_en"`
	LastnameEn  string `json:"lastname_en"`
	IsDefault   bool   `json:"is_default"`
}

type UpdateMemberBankControllerRequest = CreateMemberBankControllerRequest

func (c *Controller) CreateMemberBankController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.bank.create.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	var req CreateMemberBankControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.CreateMemberBankService(ctx.Request.Context(), memberID, &CreateMemberBankServiceRequest{
		BankID:      bankID,
		BankNo:      req.BankNo,
		FirstnameTh: req.FirstnameTh,
		LastnameTh:  req.LastnameTh,
		FirstnameEn: req.FirstnameEn,
		LastnameEn:  req.LastnameEn,
		IsDefault:   req.IsDefault,
		ActionBy:    actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.bank.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) InfoMemberBankController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.bank.info.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, bankID, ok := c.parseMemberBankURI(ctx)
	if !ok {
		return
	}

	data, err := c.svc.InfoMemberBankService(ctx.Request.Context(), memberID, bankID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.bank.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) UpdateMemberBankController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.bank.update.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, bankRowID, ok := c.parseMemberBankURI(ctx)
	if !ok {
		return
	}

	var req UpdateMemberBankControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.UpdateMemberBankService(ctx.Request.Context(), memberID, bankRowID, &UpdateMemberBankServiceRequest{
		BankID:      bankID,
		BankNo:      req.BankNo,
		FirstnameTh: req.FirstnameTh,
		LastnameTh:  req.LastnameTh,
		FirstnameEn: req.FirstnameEn,
		LastnameEn:  req.LastnameEn,
		IsDefault:   req.IsDefault,
		ActionBy:    actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.bank.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteMemberBankController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.bank.delete.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, bankID, ok := c.parseMemberBankURI(ctx)
	if !ok {
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.DeleteMemberBankService(ctx.Request.Context(), memberID, bankID, actionBy); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.bank.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseMemberBankURI(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	var uri MemberBankURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	memberID, err := uuid.Parse(uri.MemberID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	bankID, err := uuid.Parse(uri.BankID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return memberID, bankID, true
}
