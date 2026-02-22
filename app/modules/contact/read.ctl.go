package contact

import (
	"phakram/app/utils/base"

	"github.com/gin-gonic/gin"
)

func (c *Controller) MarkReadController(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.MarkRead(ctx.Request.Context(), id, true); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}

func (c *Controller) MarkUnreadController(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.MarkRead(ctx.Request.Context(), id, false); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}
