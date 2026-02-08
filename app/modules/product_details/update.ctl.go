package product_details

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateProductDetailControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateProductDetailController struct {
	ProductID        string           `json:"product_id"`
	Description      string           `json:"description"`
	Material         string           `json:"material"`
	Dimensions       string           `json:"dimensions"`
	Weight           *decimal.Decimal `json:"weight"`
	CareInstructions string           `json:"care_instructions"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateProductDetailControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_details.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateProductDetailController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_details.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var productID uuid.UUID
	if req.ProductID != "" {
		productID, err = uuid.Parse(req.ProductID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateProductDetailService{
		ProductID:        productID,
		Description:      req.Description,
		Material:         req.Material,
		Dimensions:       req.Dimensions,
		Weight:           req.Weight,
		CareInstructions: req.CareInstructions,
		MemberID:         memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`product_details.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) ProductDetailsUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
