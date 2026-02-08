package dashboard

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type SummaryRequest struct {
	Range string `form:"range"`
}

func (c *Controller) Summary(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`dashboard.ctl.summary.start`)

	if !authmod.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var req SummaryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.SummaryService(ctx.Request.Context(), req.Range)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`dashboard.ctl.summary.callsvc`)

	base.Success(ctx, data)
}
