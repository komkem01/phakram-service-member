package prefixes

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InfoPrefixControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoPrefixControllerResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	GenderID  uuid.UUID `json:"gender_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoPrefixControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`prefixes.ctl.info.request`)

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
	span.AddEvent(`prefixes.ctl.info.callsvc`)

	var resp InfoPrefixControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) PrefixesInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
