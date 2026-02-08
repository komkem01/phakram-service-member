package member_accounts

import (
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateMemberAccountController struct {
	MemberID string `json:"member_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Controller) CreateMemberAccountController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`member_accounts.ctl.create.start`)

	var req CreateMemberAccountController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_accounts.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if err := c.svc.CreateMemberAccountService(ctx.Request.Context(), &CreateMemberAccountService{
		MemberID: memberID,
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`member_accounts.ctl.create.success`)
	base.Success(ctx, nil)
}
