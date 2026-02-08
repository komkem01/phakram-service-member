package tiers

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateTierControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateTierController struct {
	NameTh       string           `json:"name_th"`
	NameEn       string           `json:"name_en"`
	MinSpending  *decimal.Decimal `json:"min_spending"`
	IsActive     *bool            `json:"is_active"`
	DiscountRate *decimal.Decimal `json:"discount_rate"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateTierControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`tiers.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateTierController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`tiers.ctl.update.request_body`)

	if err := c.svc.UpdateService(ctx, id, &UpdateTierService{
		NameTh:       req.NameTh,
		NameEn:       req.NameEn,
		MinSpending:  req.MinSpending,
		IsActive:     req.IsActive,
		DiscountRate: req.DiscountRate,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`tiers.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) TiersUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
