package tiers

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type CreateTierController struct {
	NameTh       string `json:"name_th"`
	NameEn       string `json:"name_en"`
	MinSpending  string `json:"min_spending"`
	IsActive     *bool  `json:"is_active"`
	DiscountRate string `json:"discount_rate"`
}

func (c *Controller) CreateTierController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`tiers.ctl.create.start`)

	var req CreateTierController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`tiers.ctl.create.request`)

	minSpendingDec, err := decimal.NewFromString(req.MinSpending)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	discountRateDec, err := decimal.NewFromString(req.DiscountRate)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	if err := c.svc.CreateTierService(ctx.Request.Context(), &CreateTierService{
		NameTh:       req.NameTh,
		NameEn:       req.NameEn,
		MinSpending:  minSpendingDec,
		IsActive:     req.IsActive,
		DiscountRate: discountRateDec,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`tiers.ctl.create.success`)
	base.Success(ctx, nil)
}
