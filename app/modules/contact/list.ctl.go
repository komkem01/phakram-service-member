package contact

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type ListContactControllerRequest struct {
	base.RequestPaginate
	SendStatus string `form:"send_status"`
	ReadStatus string `form:"read_status"`
}

func (c *Controller) ListController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListContactControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any("body", req)).Errf("internal: %s", err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent("contact.ctl.list.request")

	data, page, err := c.svc.List(ctx.Request.Context(), &ListContactServiceRequest{
		RequestPaginate: req.RequestPaginate,
		SendStatus:      req.SendStatus,
		ReadStatus:      req.ReadStatus,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("contact.ctl.list.success")
	base.Paginate(ctx, data, page)
}
