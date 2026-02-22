package contact

import (
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/gin-gonic/gin"
)

func (c *Controller) InfoController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	id := ctx.Param("id")

	data, err := c.svc.Info(ctx.Request.Context(), id)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("contact.ctl.info.success")
	base.Success(ctx, data)
}
