package products

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListProductControllerRequest struct {
	base.RequestPaginate
}

type ListProductControllerResponses struct {
	ID         uuid.UUID       `json:"id"`
	CategoryID uuid.UUID       `json:"category_id"`
	NameTh     string          `json:"name_th"`
	NameEn     string          `json:"name_en"`
	ProductNo  string          `json:"product_no"`
	Price      decimal.Decimal `json:"price"`
	ImageURL   string          `json:"image_url,omitempty"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

func (c *Controller) ProductsList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListProductControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`products.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListProductServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`products.ctl.list.callsvc`)

	var resp []*ListProductControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
