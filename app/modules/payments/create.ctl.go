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

type CreatePaymentController struct {
	Amount     decimal.Decimal `json:"amount"`
	Status     string          `json:"status"`
	ApprovedBy string          `json:"approved_by"`
	ApprovedAt time.Time       `json:"approved_at"`
}

func (c *Controller) CreatePaymentController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`payments.ctl.create.start`)

	var req CreatePaymentController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payments.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	approvedBy, err := uuid.Parse(req.ApprovedBy)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreatePaymentService(ctx.Request.Context(), &CreatePaymentService{
		Amount:     req.Amount,
		Status:     req.Status,
		ApprovedBy: approvedBy,
		ApprovedAt: req.ApprovedAt,
		MemberID:   memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`payments.ctl.create.success`)
	base.Success(ctx, nil)
}
