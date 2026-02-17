package categories

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateCategoryController struct {
	ParentID *string `json:"parent_id"`
	NameTh   string  `json:"name_th"`
	NameEn   string  `json:"name_en"`
	IsActive *bool   `json:"is_active"`
}

func (c *Controller) CreateCategoryController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`categories.ctl.create.start`)

	var req CreateCategoryController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`categories.ctl.create.request`)

	if err := c.svc.CreateCategoryService(ctx.Request.Context(), &CreateCategoryService{
		ParentID: req.ParentID,
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		IsActive: req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`categories.ctl.create.success`)
	base.Success(ctx, nil)
}
