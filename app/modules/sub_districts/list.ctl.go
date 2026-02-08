package sub_districts

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListSubDistrictControllerRequest struct {
	base.RequestPaginate
}

type ListSubDistrictControllerResponses struct {
	ID         uuid.UUID `json:"id"`
	DistrictID uuid.UUID `json:"district_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func (c *Controller) SubDistrictsList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListSubDistrictControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`sub_districts.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListSubDistrictServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`sub_districts.ctl.list.callsvc`)

	var resp []*ListSubDistrictControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
