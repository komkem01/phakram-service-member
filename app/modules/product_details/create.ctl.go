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

type CreateProductDetailController struct {
	ProductID        string          `json:"product_id"`
	Description      string          `json:"description"`
	Material         string          `json:"material"`
	Dimensions       string          `json:"dimensions"`
	Weight           decimal.Decimal `json:"weight"`
	CareInstructions string          `json:"care_instructions"`
}

func (c *Controller) CreateProductDetailController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`product_details.ctl.create.start`)

	var req CreateProductDetailController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_details.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateProductDetailService(ctx.Request.Context(), &CreateProductDetailService{
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

	span.AddEvent(`product_details.ctl.create.success`)
	base.Success(ctx, nil)
}
