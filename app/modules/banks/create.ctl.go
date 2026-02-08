package banks

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateBankController struct {
	NameTh    string `json:"name_th"`
	NameAbbTh string `json:"name_abb_th"`
	NameEn    string `json:"name_en"`
	NameAbbEn string `json:"name_abb_en"`
	IsActive  bool   `json:"is_active"`
}

func (c *Controller) CreateBankController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`banks.ctl.create.start`)

	var req CreateBankController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`banks.ctl.create.request`)

	if err := c.svc.CreateBankService(ctx.Request.Context(), &CreateBankService{
		NameTh:    req.NameTh,
		NameAbbTh: req.NameAbbTh,
		NameEn:    req.NameEn,
		NameAbbEn: req.NameAbbEn,
		IsActive:  req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`banks.ctl.create.success`)
	base.Success(ctx, nil)
}
