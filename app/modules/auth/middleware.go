package auth

import (
	"strings"

	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

func (c *Controller) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span, log := utils.LogSpanFromGin(ctx)
		span.AddEvent(`auth.middleware.start`)

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			ctx.Abort()
			return
		}

		claims, err := c.svc.parseToken(parts[1], "access")
		if err != nil {
			log.Errf(`internal: %s`, err)
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			ctx.Abort()
			return
		}

		ctx.Set(ContextMemberIDKey, claims.Sub)
		ctx.Set(ContextRoleKey, claims.Role)
		ctx.Set(ContextIsAdminKey, claims.IsAdmin)

		span.AddEvent(`auth.middleware.success`)
		ctx.Next()
	}
}
