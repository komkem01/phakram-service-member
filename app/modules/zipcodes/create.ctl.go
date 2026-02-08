package zipcodes

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateZipcodeController struct {
	SubDistrictsID uuid.UUID `json:"sub_districts_id"`
	Name           string    `json:"name"`
	IsActive       bool      `json:"is_active"`
}

func (c *Controller) CreateZipcodeController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`zipcodes.ctl.create.start`)

	var req CreateZipcodeController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`zipcodes.ctl.create.request`)

	if err := c.svc.CreateZipcodeService(ctx.Request.Context(), &CreateZipcodeService{
		SubDistrictsID: req.SubDistrictsID,
		Name:           req.Name,
		IsActive:       req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`zipcodes.ctl.create.success`)
	base.Success(ctx, nil)
}
