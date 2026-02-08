package products

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

type UpdateProductControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateProductController struct {
	CategoryID string           `json:"category_id"`
	NameTh     string           `json:"name_th"`
	NameEn     string           `json:"name_en"`
	ProductNo  string           `json:"product_no"`
	Price      *decimal.Decimal `json:"price"`
	IsActive   *bool            `json:"is_active"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateProductControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`products.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateProductController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`products.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var categoryID uuid.UUID
	if req.CategoryID != "" {
		categoryID, err = uuid.Parse(req.CategoryID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateProductService{
		CategoryID: categoryID,
		NameTh:     req.NameTh,
		NameEn:     req.NameEn,
		ProductNo:  req.ProductNo,
		Price:      req.Price,
		IsActive:   req.IsActive,
		MemberID:   memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`products.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) ProductsUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
