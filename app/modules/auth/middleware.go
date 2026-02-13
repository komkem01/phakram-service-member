package auth

import (
	"strings"

	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		actorSub := claims.Sub
		actorIsAdmin := claims.IsAdmin
		if claims.ActorSub != "" {
			actorSub = claims.ActorSub
			actorIsAdmin = claims.ActorIsAdmin
		}

		ctx.Set(ContextActorMemberIDKey, actorSub)
		ctx.Set(ContextActorIsAdminKey, actorIsAdmin)
		ctx.Set(ContextActingAsKey, claims.ActingAs)

		targetID, targetErr := uuid.Parse(claims.Sub)
		actorID, actorErr := uuid.Parse(actorSub)
		if targetErr == nil && actorErr == nil {
			ctx.Request = ctx.Request.WithContext(WithRequestMeta(
				ctx.Request.Context(),
				ctx.FullPath(),
				actorID,
				targetID,
				claims.ActingAs,
			))
		}

		span.AddEvent(`auth.middleware.success`)
		ctx.Next()
	}
}
