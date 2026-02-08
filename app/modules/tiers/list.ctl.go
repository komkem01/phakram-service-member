package tiers

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListTierControllerRequest struct {
	base.RequestPaginate
}

type ListTierControllerResponses struct {
	ID           uuid.UUID `json:"id"`
	NameTh       string    `json:"name_th"`
	NameEn       string    `json:"name_en"`
	MinSpending  float64   `json:"min_spending"`
	IsActive     bool      `json:"is_active"`
	DiscountRate float64   `json:"discount_rate"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func (c *Controller) TiersList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListTierControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`tiers.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListTierServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`tiers.ctl.list.callsvc`)

	var resp []*ListTierControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
