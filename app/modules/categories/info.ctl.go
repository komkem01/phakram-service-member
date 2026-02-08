package categories

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InfoCategoryControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoCategoryControllerResponses struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	NameTh    string     `json:"name_th"`
	NameEn    string     `json:"name_en"`
	IsActive  bool       `json:"is_active"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoCategoryControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`categories.ctl.info.request`)

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
	span.AddEvent(`categories.ctl.info.callsvc`)

	var resp InfoCategoryControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) CategoriesInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
