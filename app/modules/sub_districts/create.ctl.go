package sub_districts

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateSubDistrictController struct {
	DistrictID uuid.UUID `json:"district_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
}

func (c *Controller) CreateSubDistrictController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`sub_districts.ctl.create.start`)

	var req CreateSubDistrictController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`sub_districts.ctl.create.request`)

	if err := c.svc.CreateSubDistrictService(ctx.Request.Context(), &CreateSubDistrictService{
		DistrictID: req.DistrictID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`sub_districts.ctl.create.success`)
	base.Success(ctx, nil)
}
