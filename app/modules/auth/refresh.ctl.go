package auth

import (
	"errors"
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type RefreshTokenController struct {
	RefreshToken string `json:"refresh_token"`
}

func (c *Controller) RefreshTokenController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.refresh.start`)

	var req RefreshTokenController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`members.ctl.refresh.request`)

	resp, err := c.svc.RefreshTokenService(ctx.Request.Context(), &RefreshTokenService{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			return
		}
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.refresh.success`)
	base.Success(ctx, resp)
}
