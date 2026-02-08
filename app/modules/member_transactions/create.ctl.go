package member_transactions

import (
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreateMemberTransactionController struct {
	MemberID string `json:"member_id"`
	Action   string `json:"action"`
	Details  string `json:"details"`
}

func (c *Controller) CreateMemberTransactionController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`member_transactions.ctl.create.start`)

	var req CreateMemberTransactionController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_transactions.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if err := c.svc.CreateMemberTransactionService(ctx.Request.Context(), &CreateMemberTransactionService{
		MemberID: memberID,
		Action:   req.Action,
		Details:  req.Details,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`member_transactions.ctl.create.success`)
	base.Success(ctx, nil)
}
