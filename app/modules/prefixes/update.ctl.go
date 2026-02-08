package prefixes

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdatePrefixControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdatePrefixController struct {
	NameTh   string `json:"name_th"`
	NameEn   string `json:"name_en"`
	GenderID string `json:"gender_id"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdatePrefixControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`prefixes.ctl.update.request_uri`)

	var req UpdatePrefixController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.UpdateService(ctx, id, &UpdatePrefixService{
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		GenderID: genderID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`prefixes.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) PrefixesUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
