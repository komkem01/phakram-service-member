package members

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberPaymentURIRequest struct {
	MemberID  string `uri:"id"`
	PaymentID string `uri:"member_payment_id"`
}

type CreateMemberPaymentControllerRequest struct {
	PaymentID string `json:"payment_id"`
	Quantity  int    `json:"quantity"`
	Price     string `json:"price"`
}

type UpdateMemberPaymentControllerRequest = CreateMemberPaymentControllerRequest

type ListMemberPaymentsControllerRequest struct {
	base.RequestPaginate
}

func (c *Controller) ListMemberPaymentsController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.payment.list.start`)

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req ListMemberPaymentsControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListMemberPaymentsService(ctx.Request.Context(), &ListMemberPaymentsServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.payment.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) CreateMemberPaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.payment.create.start`)

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req CreateMemberPaymentControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.CreateMemberPaymentService(ctx.Request.Context(), memberID, &CreateMemberPaymentServiceRequest{
		PaymentID: paymentID,
		Quantity:  req.Quantity,
		Price:     req.Price,
		ActionBy:  actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.payment.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) InfoMemberPaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.payment.info.start`)

	memberID, rowID, ok := c.parseMemberPaymentURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	data, err := c.svc.InfoMemberPaymentService(ctx.Request.Context(), memberID, rowID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.payment.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) UpdateMemberPaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.payment.update.start`)

	memberID, rowID, ok := c.parseMemberPaymentURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req UpdateMemberPaymentControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.UpdateMemberPaymentService(ctx.Request.Context(), memberID, rowID, &UpdateMemberPaymentServiceRequest{
		PaymentID: paymentID,
		Quantity:  req.Quantity,
		Price:     req.Price,
		ActionBy:  actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.payment.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteMemberPaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.payment.delete.start`)

	memberID, rowID, ok := c.parseMemberPaymentURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.DeleteMemberPaymentService(ctx.Request.Context(), memberID, rowID, actionBy); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.payment.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseMemberPaymentURI(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	var uri MemberPaymentURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	memberID, err := uuid.Parse(uri.MemberID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	paymentID, err := uuid.Parse(uri.PaymentID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return memberID, paymentID, true
}
