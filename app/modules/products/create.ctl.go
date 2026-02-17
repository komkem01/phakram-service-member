package products

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type CreateProductController struct {
	CategoryID string `json:"category_id"`
	NameTh     string `json:"name_th"`
	NameEn     string `json:"name_en"`
	Price      string `json:"price"`
	IsActive   *bool  `json:"is_active"`
}

func (c *Controller) CreateProductController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`products.ctl.create.start`)

	var req CreateProductController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`products.ctl.create.request`)

	priceDec, err := decimal.NewFromString(req.Price)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateProductService(ctx.Request.Context(), &CreateProductService{
		CategoryID: req.CategoryID,
		NameTh:     req.NameTh,
		NameEn:     req.NameEn,
		Price:      priceDec,
		IsActive:   req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`products.ctl.create.success`)
	base.Success(ctx, nil)
}
