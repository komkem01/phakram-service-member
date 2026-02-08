package genders

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateGenderControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateGenderController struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateGenderControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`genders.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateGenderController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`genders.ctl.update.request_body`)

	if err := c.svc.UpdateService(ctx, id, &UpdateGenderService{
		NameTh: req.NameTh,
		NameEn: req.NameEn,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`genders.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) GendersUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
