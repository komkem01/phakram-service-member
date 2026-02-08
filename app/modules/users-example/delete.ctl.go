package users

// import (
// 	"log/slog"
// 	"phakram-craft/app/modules/auth"
// 	"phakram-craft/app/modules/entities/ent"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/base"
// 	"phakram-craft/config/i18n"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type DeleteControllerRequest struct {
// 	ID string `uri:"id"`
// }

// func (c *Controller) DeleteController(ctx *gin.Context) {
// 	span, log := utils.LogSpanFromGin(ctx)

// 	var req DeleteControllerRequest
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.delete.request`)

// 	id, err := uuid.Parse(req.ID)
// 	if err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}

// 	actorRole, exists := auth.GetUserRole(ctx)
// 	if !exists {
// 		base.Unauthorized(ctx, i18n.Unauthorized, nil)
// 		return
// 	}

// 	if err := c.svc.DeleteService(ctx, ent.UserRole(actorRole), id); err != nil {
// 		base.HandleError(ctx, err)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.delete.callsvc`)

// 	base.Success(ctx, nil)
// }
