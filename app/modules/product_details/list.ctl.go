package product_details

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListProductDetailControllerRequest struct {
	base.RequestPaginate
}

type ListProductDetailControllerResponses struct {
	ID               uuid.UUID       `json:"id"`
	ProductID        uuid.UUID       `json:"product_id"`
	Description      string          `json:"description"`
	Material         string          `json:"material"`
	Dimensions       string          `json:"dimensions"`
	Weight           decimal.Decimal `json:"weight"`
	CareInstructions string          `json:"care_instructions"`
}

func (c *Controller) ProductDetailsList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListProductDetailControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_details.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListProductDetailServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`product_details.ctl.list.callsvc`)

	var resp []*ListProductDetailControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
