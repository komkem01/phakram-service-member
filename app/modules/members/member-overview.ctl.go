package members

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type MemberOverviewControllerRequest struct {
	base.RequestPaginate
}

func (c *Controller) MemberOverviewController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.overview.start`)

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req MemberOverviewControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.MemberOverviewService(ctx.Request.Context(), &MemberOverviewServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.overview.success`)
	base.Success(ctx, data)
}
