package systembankaccounts

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListSystemBankAccountControllerRequest struct {
	base.RequestPaginate
	BankID   string `form:"bank_id"`
	IsActive *bool  `form:"is_active"`
}

type UpsertSystemBankAccountControllerRequest struct {
	BankID           string `json:"bank_id"`
	AccountName      string `json:"account_name"`
	AccountNo        string `json:"account_no"`
	Branch           string `json:"branch"`
	IsActive         bool   `json:"is_active"`
	IsDefaultReceive bool   `json:"is_default_receive"`
	IsDefaultRefund  bool   `json:"is_default_refund"`
}

type InfoSystemBankAccountControllerRequest struct {
	ID string `uri:"id"`
}

type ListSystemBankAccountControllerResponse struct {
	ID               uuid.UUID `json:"id"`
	BankID           uuid.UUID `json:"bank_id"`
	BankNameTh       string    `json:"bank_name_th"`
	BankNameEn       string    `json:"bank_name_en"`
	AccountName      string    `json:"account_name"`
	AccountNo        string    `json:"account_no"`
	Branch           string    `json:"branch"`
	IsActive         bool      `json:"is_active"`
	IsDefaultReceive bool      `json:"is_default_receive"`
	IsDefaultRefund  bool      `json:"is_default_refund"`
	CreatedAt        int64     `json:"created_at"`
	UpdatedAt        int64     `json:"updated_at"`
}

func (c *Controller) ListController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListSystemBankAccountControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`system_bank_accounts.ctl.list.request`)

	bankID := uuid.Nil
	if req.BankID != "" {
		parsedBankID, err := uuid.Parse(req.BankID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		bankID = parsedBankID
	}

	data, page, err := c.svc.ListService(ctx, &ListSystemBankAccountServiceRequest{
		RequestPaginate: req.RequestPaginate,
		BankID:          bankID,
		IsActive:        req.IsActive,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	var resp []*ListSystemBankAccountControllerResponse
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}

func (c *Controller) InfoController(ctx *gin.Context) {
	var req InfoSystemBankAccountControllerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.InfoService(ctx.Request.Context(), req.ID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	resp := new(ListSystemBankAccountControllerResponse)
	if err := utils.CopyNTimeToUnix(resp, data); err != nil {
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) CreateController(ctx *gin.Context) {
	var req UpsertSystemBankAccountControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateService(ctx.Request.Context(), &UpsertSystemBankAccountServiceRequest{
		BankID:           bankID,
		AccountName:      req.AccountName,
		AccountNo:        req.AccountNo,
		Branch:           req.Branch,
		IsActive:         req.IsActive,
		IsDefaultReceive: req.IsDefaultReceive,
		IsDefaultRefund:  req.IsDefaultRefund,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	var uriReq InfoSystemBankAccountControllerRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpsertSystemBankAccountControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.UpdateService(ctx.Request.Context(), uriReq.ID, &UpsertSystemBankAccountServiceRequest{
		BankID:           bankID,
		AccountName:      req.AccountName,
		AccountNo:        req.AccountNo,
		Branch:           req.Branch,
		IsActive:         req.IsActive,
		IsDefaultReceive: req.IsDefaultReceive,
		IsDefaultRefund:  req.IsDefaultRefund,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	var req InfoSystemBankAccountControllerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.DeleteService(ctx.Request.Context(), req.ID); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}
