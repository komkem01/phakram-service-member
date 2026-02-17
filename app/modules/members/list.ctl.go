package members

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type ListControllerRequest struct {
	base.RequestPaginate
	Role string `json:"role" form:"role"`
}

func (c *Controller) ListController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.list.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var req ListControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListService(ctx.Request.Context(), &ListServiceRequest{
		RequestPaginate: req.RequestPaginate,
		Role:            req.Role,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.list.success`)
	base.Paginate(ctx, data, page)
}
