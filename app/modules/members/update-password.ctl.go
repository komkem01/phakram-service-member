package members

import (
	"strings"

	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdatePasswordControllerRequestURI struct {
	ID string `uri:"id"`
}

type UpdatePasswordControllerRequest struct {
	Password string `json:"password"`
}

func (c *Controller) UpdatePasswordController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.update_password.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var uri UpdatePasswordControllerRequestURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdatePasswordControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var actionBy *uuid.UUID
	if memberID, ok := auth.GetMemberID(ctx); ok {
		actionBy = &memberID
	}

	if err := c.svc.UpdatePasswordService(ctx.Request.Context(), id, &UpdatePasswordServiceRequest{
		Password: req.Password,
		ActionBy: actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.update_password.success`)
	base.Success(ctx, nil)
}
