package payments

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type UpdatePaymentController struct {
	Amount string `json:"amount"`
	Status string `json:"status"`
}

func (c *Controller) PaymentsUpdate(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`payments.ctl.update.start`)

	var reqUri InfoPaymentControllerRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var reqBody UpdatePaymentController
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.UpdatePaymentService(ctx.Request.Context(), reqUri.ID, &UpdatePaymentService{Amount: reqBody.Amount, Status: reqBody.Status}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`payments.ctl.update.success`)
	base.Success(ctx, nil)
}
