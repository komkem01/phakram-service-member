package users

// import (
// 	"log/slog"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/base"
// 	"phakram-craft/config/i18n"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type InfoUserControllerRequestUri struct {
// 	ID string `uri:"id"`
// }

// type InfoUserController struct {
// 	ID        string `json:"id"`
// 	Username  string `json:"username"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Email     string `json:"email"`
// 	Phone     string `json:"phone"`
// 	Status    string `json:"status"`
// 	Role      string `json:"role"`
// 	CreatedAt string `json:"created_at"`
// 	UpdatedAt string `json:"updated_at"`
// }

// func (c *Controller) InfoController(ctx *gin.Context) {
// 	span, log := utils.LogSpanFromGin(ctx)

// 	var req InfoUserControllerRequestUri
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.info.request`)

// 	id, err := uuid.Parse(req.ID)
// 	if err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	data, err := c.svc.InfoService(ctx, id)
// 	if err != nil {
// 		base.HandleError(ctx, err)
// 		return
// 	}
// 	var resp InfoUserController
// 	span.AddEvent(`users.ctl.info.callsvc`)
// 	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.InternalServerError(ctx, err.Error(), nil)
// 		return
// 	}

// 	base.Success(ctx, resp)
// }
