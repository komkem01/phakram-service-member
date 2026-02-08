package districts

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateDistrictControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateDistrictController struct {
	ProvinceID *uuid.UUID `json:"province_id"`
	Name       string     `json:"name"`
	IsActive   *bool      `json:"is_active"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateDistrictControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`districts.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateDistrictController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`districts.ctl.update.request_body`)

	if err := c.svc.UpdateService(ctx, id, &UpdateDistrictService{
		ProvinceID: req.ProvinceID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`districts.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) DistrictsUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
