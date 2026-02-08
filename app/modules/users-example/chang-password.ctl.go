package users

// import (
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/base"
// 	"phakram-craft/config/i18n"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type ChangePasswordController struct {
// 	OldPassword string `json:"old_password" binding:"required,min=6"`
// 	NewPassword string `json:"new_password" binding:"required,min=6"`
// }

// func (c *Controller) ChangePasswordController(ctx *gin.Context) {
// 	span, log := utils.LogSpanFromGin(ctx)
// 	span.AddEvent(`users.ctl.change_password.start`)

// 	// Get user ID from URI
// 	userIDStr := ctx.Param("id")
// 	userID, err := uuid.Parse(userIDStr)
// 	if err != nil {
// 		log.Errf("Invalid user ID: %v", err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.change_password.user_id_parsed`)

// 	// Bind request
// 	var req ChangePasswordController
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.Errf("Failed to bind request: %v", err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.change_password.request_bound`)

// 	// Validate passwords are different
// 	if req.OldPassword == req.NewPassword {
// 		log.Warnf("New password same as old password")
// 		base.BadRequest(ctx, i18n.NewPasswordSameAsOld, nil)
// 		return
// 	}

// 	// Call service
// 	err = c.svc.ChangePasswordService(ctx.Request.Context(), userID, &ChangePasswordService{
// 		OldPassword: req.OldPassword,
// 		NewPassword: req.NewPassword,
// 	})
// 	if err != nil {
// 		base.HandleError(ctx, err)
// 		return
// 	}

// 	span.AddEvent(`users.ctl.change_password.success`)
// 	base.Success(ctx, nil)
// }
