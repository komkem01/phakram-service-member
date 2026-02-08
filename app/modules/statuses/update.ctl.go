package statuses

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateStatusControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateStatusController struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateStatusControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`statuses.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateStatusController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`statuses.ctl.update.request_body`)

	if err := c.svc.UpdateService(ctx, id, &UpdateStatusService{
		NameTh: req.NameTh,
		NameEn: req.NameEn,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`statuses.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) StatusesUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
