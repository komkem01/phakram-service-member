package member_files

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListMemberFileControllerRequest struct {
	base.RequestPaginate
}

type ListMemberFileControllerResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (c *Controller) MemberFilesList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListMemberFileControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_files.ctl.list.request`)

	memberID := uuid.Nil
	isAdmin := authmod.GetIsAdmin(ctx)
	if !isAdmin {
		var ok bool
		memberID, ok = authmod.GetMemberID(ctx)
		if !ok {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			return
		}
	}

	data, page, err := c.svc.ListService(ctx, &ListMemberFileServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_files.ctl.list.callsvc`)

	var resp []*ListMemberFileControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
