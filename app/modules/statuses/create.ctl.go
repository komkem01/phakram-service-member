package statuses

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateStatusController struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
}

func (c *Controller) CreateStatusController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`statuses.ctl.create.start`)

	var req CreateStatusController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`statuses.ctl.create.request`)

	if err := c.svc.CreateStatusService(ctx.Request.Context(), &CreateStatusService{
		NameTh: req.NameTh,
		NameEn: req.NameEn,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`statuses.ctl.create.success`)
	base.Success(ctx, nil)
}
