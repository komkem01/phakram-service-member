package zipcodes

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListZipcodeControllerRequest struct {
	base.RequestPaginate
}

type ListZipcodeControllerResponses struct {
	ID             uuid.UUID `json:"id"`
	SubDistrictsID uuid.UUID `json:"sub_districts_id"`
	Name           string    `json:"name"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

func (c *Controller) ZipcodesList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListZipcodeControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`zipcodes.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListZipcodeServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`zipcodes.ctl.list.callsvc`)

	var resp []*ListZipcodeControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
