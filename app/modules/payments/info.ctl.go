package payments

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type InfoPaymentControllerRequest struct {
	ID string `uri:"id"`
}

func (c *Controller) PaymentsInfo(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`payments.ctl.info.start`)

	var req InfoPaymentControllerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.InfoService(ctx.Request.Context(), req.ID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`payments.ctl.info.success`)
	base.Success(ctx, data)
}
