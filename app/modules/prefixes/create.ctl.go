package prefixes

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePrefixController struct {
	NameTh   string `json:"name_th"`
	NameEn   string `json:"name_en"`
	GenderID string `json:"gender_id"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) CreatePrefixController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`prefixes.ctl.create.start`)

	var req CreatePrefixController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`prefixes.ctl.create.request`)

	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreatePrefixService(ctx.Request.Context(), &CreatePrefixService{
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		GenderID: genderID,
		IsActive: req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`prefixes.ctl.create.success`)
	base.Success(ctx, nil)
}
