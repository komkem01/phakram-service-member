package auth

import (
	"errors"
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginMemberController struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Controller) LoginMemberController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.login.start`)

	var req LoginMemberController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`members.ctl.login.request`)

	resp, err := c.svc.LoginMemberService(ctx.Request.Context(), &LoginMemberService{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			return
		}
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.login.success`)
	base.Success(ctx, resp)
}

type MeResponse struct {
	MemberID string `json:"member_id"`
}

func (c *Controller) MeController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.me.start`)

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	claims, err := c.svc.parseToken(parts[1], "access")
	if err != nil {
		log.Errf(`internal: %s`, err)
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	span.AddEvent(`members.ctl.me.success`)
	base.Success(ctx, &MeResponse{
		MemberID: claims.Sub,
	})
}
