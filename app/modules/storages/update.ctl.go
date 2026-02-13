package storages

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateStorageControllerRequest struct {
	IsActive bool `json:"is_active"`
}

func (c *Controller) StoragesUpdate(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`storages.ctl.update.start`)

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateStorageControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.UpdateService(ctx.Request.Context(), id, req.IsActive); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`storages.ctl.update.success`)
	base.Success(ctx, nil)
}
