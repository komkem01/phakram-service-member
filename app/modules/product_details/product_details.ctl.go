package productdetails

import (
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductDetailRequestUri struct {
	ID string `uri:"id"`
}

type ProductDetailRequest struct {
	Description      string `json:"description"`
	Material         string `json:"material"`
	Dimensions       string `json:"dimensions"`
	Weight           string `json:"weight"`
	CareInstructions string `json:"care_instructions"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	var req ProductDetailRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	data, err := c.svc.GetByProductID(ctx, productID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`product_details.ctl.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) CreateController(ctx *gin.Context) {
	var reqUri ProductDetailRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req ProductDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	weight, err := decimal.NewFromString(req.Weight)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	payload := &ent.ProductDetailEntity{
		Description:      req.Description,
		Material:         req.Material,
		Dimensions:       req.Dimensions,
		Weight:           weight,
		CareInstructions: req.CareInstructions,
	}
	if err := c.svc.CreateByProductID(ctx, productID, payload); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	var reqUri ProductDetailRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req ProductDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	weight, err := decimal.NewFromString(req.Weight)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	payload := &ent.ProductDetailEntity{
		Description:      req.Description,
		Material:         req.Material,
		Dimensions:       req.Dimensions,
		Weight:           weight,
		CareInstructions: req.CareInstructions,
	}
	if err := c.svc.UpdateByProductID(ctx, productID, payload); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	var reqUri ProductDetailRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	if err := c.svc.DeleteByProductID(ctx, productID); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}
