package payments

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type CreatePaymentController struct {
	Amount string `json:"amount"`
	Status string `json:"status"`
}

func (c *Controller) CreatePaymentController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`payments.ctl.create.start`)

	var req CreatePaymentController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreatePaymentService(ctx.Request.Context(), &CreatePaymentService{Amount: req.Amount, Status: req.Status}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`payments.ctl.create.success`)
	base.Success(ctx, nil)
}
