package zipcodes

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateZipcodeControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateZipcodeController struct {
	SubDistrictsID *uuid.UUID `json:"sub_districts_id"`
	Name           string     `json:"name"`
	IsActive       *bool      `json:"is_active"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateZipcodeControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`zipcodes.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateZipcodeController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`zipcodes.ctl.update.request_body`)

	if err := c.svc.UpdateService(ctx, id, &UpdateZipcodeService{
		SubDistrictsID: req.SubDistrictsID,
		Name:           req.Name,
		IsActive:       req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`zipcodes.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) ZipcodesUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
