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

type UpdatePaymentFileControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdatePaymentFileController struct {
	PaymentID string `json:"payment_id"`
	FileID    string `json:"file_id"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdatePaymentFileControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payment_files.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdatePaymentFileController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payment_files.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var paymentID uuid.UUID
	if req.PaymentID != "" {
		paymentID, err = uuid.Parse(req.PaymentID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var fileID uuid.UUID
	if req.FileID != "" {
		fileID, err = uuid.Parse(req.FileID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdatePaymentFileService{
		PaymentID: paymentID,
		FileID:    fileID,
		MemberID:  memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`payment_files.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) PaymentFilesUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
