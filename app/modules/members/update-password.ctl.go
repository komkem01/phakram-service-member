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
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	Password        string `json:"password"`
}

func (c *Controller) UpdatePasswordController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.update_password.start`)

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

	if !c.ensureAdminOrSelf(ctx, id) {
		return
	}

	var req UpdatePasswordControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	password := strings.TrimSpace(req.NewPassword)
	if password == "" {
		password = strings.TrimSpace(req.Password)
	}

	if password == "" {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var actionBy *uuid.UUID
	if memberID, ok := auth.GetMemberID(ctx); ok {
		actionBy = &memberID
	}

	if err := c.svc.UpdatePasswordService(ctx.Request.Context(), id, &UpdatePasswordServiceRequest{
		Password: password,
		ActionBy: actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.update_password.success`)
	base.Success(ctx, nil)
}
