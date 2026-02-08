package payments

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdatePaymentControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdatePaymentController struct {
	Amount     *decimal.Decimal `json:"amount"`
	Status     string           `json:"status"`
	ApprovedBy string           `json:"approved_by"`
	ApprovedAt *time.Time       `json:"approved_at"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdatePaymentControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payments.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdatePaymentController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payments.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var approvedBy uuid.UUID
	if req.ApprovedBy != "" {
		approvedBy, err = uuid.Parse(req.ApprovedBy)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdatePaymentService{
		Amount:     req.Amount,
		Status:     req.Status,
		ApprovedBy: approvedBy,
		ApprovedAt: req.ApprovedAt,
		MemberID:   memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`payments.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) PaymentsUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
