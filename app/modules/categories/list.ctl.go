package categories

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListCategoryControllerRequest struct {
	base.RequestPaginate
}

type ListCategoryControllerResponses struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	NameTh    string     `json:"name_th"`
	NameEn    string     `json:"name_en"`
	IsActive  bool       `json:"is_active"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

func (c *Controller) CategoriesList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListCategoryControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`categories.ctl.list.request`)

	data, page, err := c.svc.ListService(ctx, &ListCategoryServiceRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`categories.ctl.list.callsvc`)

	var resp []*ListCategoryControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
