package users

// import (
// 	"log/slog"
// 	"phakram-craft/app/modules/entities/ent"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/base"
// 	"phakram-craft/config/i18n"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type ListUserControllerRequest struct {
// 	base.RequestPaginate
// }

// type ListUserControllerResponses struct {
// 	ID        uuid.UUID `json:"id"`
// 	Username  string    `json:"username"`
// 	FirstName string    `json:"first_name"`
// 	LastName  string    `json:"last_name"`
// 	Status    string    `json:"status"`
// 	Role      string    `json:"role"`
// 	Email     string    `json:"email"`
// 	Phone     string    `json:"phone"`
// 	CreatedAt string    `json:"created_at"`
// 	UpdatedAt string    `json:"updated_at"`
// }

// func (c *Controller) ListController(ctx *gin.Context) {
// 	span, log := utils.LogSpanFromGin(ctx)

// 	var req ListUserControllerRequest
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.list.request`)

// 	data, page, err := c.svc.ListService(ctx, &ListServiceRequest{
// 		RequestPaginate: req.RequestPaginate,
// 	})
// 	if err != nil {
// 		base.HandleError(ctx, err)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.list.callsvc`)

// 	var resp []*ListUserControllerResponses
// 	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.InternalServerError(ctx, err.Error(), nil)
// 		return
// 	}

// 	base.Paginate(ctx, resp, page)
// }

// func (c *Controller) ListCustomerController(ctx *gin.Context) {
// 	span, log := utils.LogSpanFromGin(ctx)

// 	var req ListUserControllerRequest
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.list_customer.request`)

// 	data, page, err := c.svc.ListService(ctx, &ListServiceRequest{
// 		RequestPaginate: req.RequestPaginate,
// 		Role:            string(ent.UserRoleUser),
// 	})
// 	if err != nil {
// 		base.HandleError(ctx, err)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.list_customer.callsvc`)

// 	var resp []*ListUserControllerResponses
// 	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		base.InternalServerError(ctx, err.Error(), nil)
// 		return
// 	}

// 	base.Paginate(ctx, resp, page)
// }
