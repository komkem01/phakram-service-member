package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextMemberIDKey = "member_id"
	ContextRoleKey     = "role"
	ContextIsAdminKey  = "is_admin"
)

func GetMemberID(ctx *gin.Context) (uuid.UUID, bool) {
	if ctx == nil {
		return uuid.Nil, false
	}

	value, ok := ctx.Get(ContextMemberIDKey)
	if !ok {
		return uuid.Nil, false
	}

	switch v := value.(type) {
	case uuid.UUID:
		return v, true
	case string:
		memberID, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, false
		}
		return memberID, true
	default:
		return uuid.Nil, false
	}
}

func GetIsAdmin(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}

	value, ok := ctx.Get(ContextIsAdminKey)
	if !ok {
		return false
	}

	isAdmin, ok := value.(bool)
	if !ok {
		return false
	}

	return isAdmin
}
