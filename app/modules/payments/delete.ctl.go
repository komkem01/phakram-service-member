package payments

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

func (c *Controller) PaymentsDelete(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`payments.ctl.delete.start`)

	var req InfoPaymentControllerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.DeletePaymentService(ctx.Request.Context(), req.ID); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`payments.ctl.delete.success`)
	base.Success(ctx, nil)
}
