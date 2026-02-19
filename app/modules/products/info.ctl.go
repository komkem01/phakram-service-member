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

type InfoProductControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoProductControllerResponses struct {
	ID         uuid.UUID       `json:"id"`
	CategoryID uuid.UUID       `json:"category_id"`
	NameTh     string          `json:"name_th"`
	NameEn     string          `json:"name_en"`
	ProductNo  string          `json:"product_no"`
	Price      decimal.Decimal `json:"price"`
	ImageURL   string          `json:"image_url,omitempty"`
	ImageURLs  []string        `json:"image_urls,omitempty"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoProductControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`products.ctl.info.request`)

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
	span.AddEvent(`products.ctl.info.callsvc`)

	var resp InfoProductControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) ProductsInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
