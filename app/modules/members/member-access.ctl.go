package members

import (
	"phakram/app/modules/auth"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) ensureAdminOrSelf(ctx *gin.Context, targetMemberID uuid.UUID) bool {
	if auth.GetIsAdmin(ctx) {
		return true
	}

	memberID, ok := auth.GetMemberID(ctx)
	if !ok || memberID != targetMemberID {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return false
	}

	return true
}
