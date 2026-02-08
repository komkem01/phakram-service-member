package provinces

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateProvinceController struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) CreateProvinceController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`provinces.ctl.create.start`)

	var req CreateProvinceController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`provinces.ctl.create.request`)

	if err := c.svc.CreateProvinceService(ctx.Request.Context(), &CreateProvinceService{
		Name:     req.Name,
		IsActive: req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`provinces.ctl.create.success`)
	base.Success(ctx, nil)
}
