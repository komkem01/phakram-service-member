package payments

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListPaymentControllerRequest struct {
	base.RequestPaginate
	Status string `form:"status"`
}

type ListPaymentControllerResponses struct {
	ID         uuid.UUID       `json:"id"`
	Amount     decimal.Decimal `json:"amount"`
	Status     string          `json:"status"`
	ApprovedBy *uuid.UUID      `json:"approved_by"`
	ApprovedAt *string         `json:"approved_at"`
}

func (c *Controller) PaymentsList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListPaymentControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payments.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListPaymentServiceRequest{RequestPaginate: req.RequestPaginate, Status: req.Status})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`payments.ctl.list.callsvc`)

	var resp []*ListPaymentControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
