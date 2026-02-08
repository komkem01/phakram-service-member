package districts

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateDistrictController struct {
	ProvinceID uuid.UUID `json:"province_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
}

func (c *Controller) CreateDistrictController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`districts.ctl.create.start`)

	var req CreateDistrictController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`districts.ctl.create.request`)

	if err := c.svc.CreateDistrictService(ctx.Request.Context(), &CreateDistrictService{
		ProvinceID: req.ProvinceID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`districts.ctl.create.success`)
	base.Success(ctx, nil)
}
