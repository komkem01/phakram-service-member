package users

// import (
// 	"log/slog"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/base"
// 	"phakram-craft/config/i18n"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type UpdateUserControllerRequestUri struct {
// 	ID string `uri:"id"`
// }

// type UpdateUserController struct {
// 	Username  string `json:"username"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// }

// func (c *Controller) UpdateController(ctx *gin.Context) {
// 	span, log := utils.LogSpanFromGin(ctx)

// 	var reqUri UpdateUserControllerRequestUri
// 	if err := ctx.ShouldBindUri(&reqUri); err != nil {
// 		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.update.request_uri`)

// 	id, err := uuid.Parse(reqUri.ID)
// 	if err != nil {
// 		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	var req UpdateUserController
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.update.request_body`)

// 	err = c.svc.UpdateUserService(ctx, id, &UpdateUserService{
// 		Username:  req.Username,
// 		FirstName: req.FirstName,
// 		LastName:  req.LastName,
// 	})
// 	if err != nil {
// 		base.HandleError(ctx, err)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.update.callsvc`)

// 	base.Success(ctx, nil)
// }
