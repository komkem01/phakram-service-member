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

type InfoTierControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoTierControllerResponses struct {
	ID           uuid.UUID       `json:"id"`
	NameTh       string          `json:"name_th"`
	NameEn       string          `json:"name_en"`
	MinSpending  decimal.Decimal `json:"min_spending"`
	IsActive     bool            `json:"is_active"`
	DiscountRate decimal.Decimal `json:"discount_rate"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoTierControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`tiers.ctl.info.request`)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.InfoService(ctx, id)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`tiers.ctl.info.callsvc`)

	resp := &InfoTierControllerResponses{
		ID:           data.ID,
		NameTh:       data.NameTh,
		NameEn:       data.NameEn,
		MinSpending:  data.MinSpending,
		IsActive:     data.IsActive,
		DiscountRate: data.DiscountRate,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}

	base.Success(ctx, resp)
}

func (c *Controller) TiersInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
