package storages

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) StoragesInfo(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`storages.ctl.info.start`)

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.InfoService(ctx.Request.Context(), id)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`storages.ctl.info.success`)
	base.Success(ctx, data)
}
