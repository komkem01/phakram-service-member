package prefixes

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListPrefixControllerRequest struct {
	base.RequestPaginate
}

type ListPrefixControllerResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	GenderID  uuid.UUID `json:"gender_id"`
	CreatedAt string    `json:"created_at"`
}

func (c *Controller) PrefixesList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListPrefixControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`prefixes.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListPrefixServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`prefixes.ctl.list.callsvc`)

	var resp []*ListPrefixControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
