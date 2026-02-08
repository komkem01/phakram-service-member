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

type InfoProductDetailControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoProductDetailControllerResponses struct {
	ID               uuid.UUID              `json:"id"`
	ProductID        uuid.UUID              `json:"product_id"`
	Description      string                 `json:"description"`
	Material         string                 `json:"material"`
	Dimensions       string                 `json:"dimensions"`
	Weight           decimal.Decimal        `json:"weight"`
	CareInstructions string                 `json:"care_instructions"`
	Image            *InfoProductDetailFile `json:"image"`
}

type InfoProductDetailFile struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	FileID    uuid.UUID `json:"file_id"`
	FilePath  string    `json:"file_path"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoProductDetailControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_details.ctl.info.request`)

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
	span.AddEvent(`product_details.ctl.info.callsvc`)

	var resp InfoProductDetailControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) ProductDetailsInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
