package carts

import (
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateCartController struct {
	MemberID string `json:"member_id"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) CreateCartController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.create.start`)

	var req CreateCartController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`carts.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if err := c.svc.CreateCartService(ctx.Request.Context(), &CreateCartService{
		MemberID: memberID,
		IsActive: req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.create.success`)
	base.Success(ctx, nil)
}
