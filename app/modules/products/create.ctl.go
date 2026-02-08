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

type CreateProductController struct {
	CategoryID string          `json:"category_id"`
	NameTh     string          `json:"name_th"`
	NameEn     string          `json:"name_en"`
	ProductNo  string          `json:"product_no"`
	Price      decimal.Decimal `json:"price"`
	IsActive   bool            `json:"is_active"`
}

func (c *Controller) CreateProductController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`products.ctl.create.start`)

	var req CreateProductController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`products.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateProductService(ctx.Request.Context(), &CreateProductService{
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

	span.AddEvent(`products.ctl.create.success`)
	base.Success(ctx, nil)
}
