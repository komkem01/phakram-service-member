package contact

import (
	"phakram/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CountUnreadControllerResponse struct {
	Unread int64 `json:"unread"`
}

func (c *Controller) CountUnreadController(ctx *gin.Context) {
	count, err := c.svc.CountUnread(ctx.Request.Context())
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, &CountUnreadControllerResponse{
		Unread: count,
	})
}
