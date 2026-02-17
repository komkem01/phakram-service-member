package orders

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type ConfirmOrderPaymentControllerRequest struct {
	TransferredAmount string `json:"transferred_amount"`
	SlipImageBase64   string `json:"slip_image_base64"`
	SlipFileName      string `json:"slip_file_name"`
	SlipFileType      string `json:"slip_file_type"`
	SlipFileSize      int64  `json:"slip_file_size"`
}

type RejectOrderPaymentControllerRequest struct {
	Reason string `json:"reason"`
}

func (c *Controller) ConfirmOrderPaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.payment.confirm.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	var req ConfirmOrderPaymentControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.ConfirmOrderPaymentService(ctx.Request.Context(), orderID, &ConfirmOrderPaymentServiceRequest{
		TransferredAmount: req.TransferredAmount,
		SlipImageBase64:   req.SlipImageBase64,
		SlipFileName:      req.SlipFileName,
		SlipFileType:      req.SlipFileType,
		SlipFileSize:      req.SlipFileSize,
	}, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.payment.confirm.success`)
	base.Success(ctx, data)
}

func (c *Controller) ApproveOrderPaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.payment.approve.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin || !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.ApproveOrderPaymentService(ctx.Request.Context(), orderID, requesterID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.payment.approve.success`)
	base.Success(ctx, data)
}

func (c *Controller) RejectOrderPaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.payment.reject.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	var req RejectOrderPaymentControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin || !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.RejectOrderPaymentService(ctx.Request.Context(), orderID, &RejectOrderPaymentServiceRequest{Reason: req.Reason}, requesterID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.payment.reject.success`)
	base.Success(ctx, data)
}
