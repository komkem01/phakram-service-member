package payment_files

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePaymentFileController struct {
	PaymentID string `json:"payment_id"`
	FileID    string `json:"file_id"`
}

func (c *Controller) CreatePaymentFileController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`payment_files.ctl.create.start`)

	var req CreatePaymentFileController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payment_files.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	fileID, err := uuid.Parse(req.FileID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreatePaymentFileService(ctx.Request.Context(), &CreatePaymentFileService{
		PaymentID: paymentID,
		FileID:    fileID,
		MemberID:  memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`payment_files.ctl.create.success`)
	base.Success(ctx, nil)
}
