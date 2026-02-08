package genders

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateGenderController struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
}

func (c *Controller) CreateGenderController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`genders.ctl.create.start`)

	var req CreateGenderController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`genders.ctl.create.request`)

	err := c.svc.CreateGenderService(ctx.Request.Context(), &CreateGenderService{
		NameTh: req.NameTh,
		NameEn: req.NameEn,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`genders.ctl.create.success`)
	base.Success(ctx, nil)
}
