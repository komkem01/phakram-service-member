package users

// import (
// 	"phakram-craft/app/modules/auth"
// 	"phakram-craft/app/modules/entities/ent"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/base"
// 	"phakram-craft/config/i18n"

// 	"github.com/gin-gonic/gin"
// )

// type CreateUserController struct {
// 	Username  string `json:"username"`
// 	Password  string `json:"password"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Email     string `json:"email"`
// 	Phone     string `json:"phone"`
// 	Role      string `json:"role"`
// }

// func (c *Controller) CreateUserController(ctx *gin.Context) {
// 	span, _ := utils.LogSpanFromGin(ctx)
// 	span.AddEvent(`users.ctl.create.start`)

// 	var req CreateUserController
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		base.BadRequest(ctx, i18n.BadRequest, nil)
// 		return
// 	}
// 	span.AddEvent(`users.ctl.create.request`)

// 	actorRole, exists := auth.GetUserRole(ctx)
// 	if !exists {
// 		base.Unauthorized(ctx, i18n.Unauthorized, nil)
// 		return
// 	}

// 	requestedRole := ent.UserRoleUser
// 	if req.Role != "" {
// 		switch req.Role {
// 		case string(ent.UserRoleUser), string(ent.UserRoleAdmin):
// 			requestedRole = ent.UserRole(req.Role)
// 		default:
// 			base.BadRequest(ctx, i18n.BadRequest, nil)
// 			return
// 		}
// 	}

// 	switch actorRole {
// 	case string(ent.UserRoleSuperadmin):
// 		// superadmin can create admin and user only
// 	case string(ent.UserRoleAdmin):
// 		if requestedRole != ent.UserRoleUser {
// 			base.Forbidden(ctx, i18n.Forbidden, nil)
// 			return
// 		}
// 	default:
// 		base.Forbidden(ctx, i18n.Forbidden, nil)
// 		return
// 	}

// 	email := req.Email
// 	if email != "" {
// 		formattedEmail, err := utils.FormatEmail(email)
// 		if err != nil {
// 			base.BadRequest(ctx, i18n.InvalidEmailFormat, nil)
// 			return
// 		}
// 		req.Email = formattedEmail
// 		span.AddEvent(`users.ctl.create.email_formatted`)
// 	}
// 	phone := req.Phone
// 	if phone != "" {
// 		formattedPhone, err := utils.FormatPhone(phone)
// 		if err != nil {
// 			base.BadRequest(ctx, i18n.InvalidPhoneFormat, nil)
// 			return
// 		}
// 		req.Phone = formattedPhone
// 		span.AddEvent(`users.ctl.create.phone_formatted`)
// 	}

// 	err := c.svc.CreateUserService(ctx.Request.Context(), &CreateUserService{
// 		Username:  req.Username,
// 		Password:  req.Password,
// 		FirstName: req.FirstName,
// 		LastName:  req.LastName,
// 		Email:     req.Email,
// 		Phone:     req.Phone,
// 		Role:      requestedRole,
// 	})
// 	if err != nil {
// 		base.HandleError(ctx, err)
// 		return
// 	}

// 	span.AddEvent(`users.ctl.create.success`)
// 	base.Success(ctx, nil)
// }
