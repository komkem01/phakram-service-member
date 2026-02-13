package storages

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListStorageControllerRequest struct {
	RefID string `form:"ref_id"`
}

func (c *Controller) StoragesList(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`storages.ctl.list.start`)

	var req ListStorageControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	refID, err := uuid.Parse(req.RefID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.ListService(ctx.Request.Context(), refID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`storages.ctl.list.success`)
	base.Success(ctx, data)
}
