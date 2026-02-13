package storages

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) StoragesDelete(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`storages.ctl.delete.start`)

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.DeleteService(ctx.Request.Context(), id); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`storages.ctl.delete.success`)
	base.Success(ctx, nil)
}
