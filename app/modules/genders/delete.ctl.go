package genders

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteGenderControllerRequest struct {
	ID string `uri:"id"`
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req DeleteGenderControllerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`genders.ctl.delete.request`)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.DeleteService(ctx, id); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`genders.ctl.delete.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) GendersDelete(ctx *gin.Context) {
	c.DeleteController(ctx)
}
