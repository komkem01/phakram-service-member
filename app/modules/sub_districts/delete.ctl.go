package sub_districts

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteSubDistrictControllerRequest struct {
	ID string `uri:"id"`
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req DeleteSubDistrictControllerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`sub_districts.ctl.delete.request`)

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
	span.AddEvent(`sub_districts.ctl.delete.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) SubDistrictsDelete(ctx *gin.Context) {
	c.DeleteController(ctx)
}
