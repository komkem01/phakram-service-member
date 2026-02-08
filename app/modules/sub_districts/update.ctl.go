package sub_districts

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateSubDistrictControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateSubDistrictController struct {
	DistrictID *uuid.UUID `json:"district_id"`
	Name       string     `json:"name"`
	IsActive   *bool      `json:"is_active"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateSubDistrictControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`sub_districts.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateSubDistrictController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`sub_districts.ctl.update.request_body`)

	if err := c.svc.UpdateService(ctx, id, &UpdateSubDistrictService{
		DistrictID: req.DistrictID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`sub_districts.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) SubDistrictsUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
