package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextMemberIDKey      = "member_id"
	ContextRoleKey          = "role"
	ContextIsAdminKey       = "is_admin"
	ContextActorMemberIDKey = "actor_member_id"
	ContextActorIsAdminKey  = "actor_is_admin"
	ContextActingAsKey      = "is_acting_as"
)

type requestContextKey string

const (
	requestEndpointKey     requestContextKey = "request_endpoint"
	requestActorMemberKey  requestContextKey = "request_actor_member_id"
	requestTargetMemberKey requestContextKey = "request_target_member_id"
	requestActingAsKey     requestContextKey = "request_is_acting_as"
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

func GetActorMemberID(ctx *gin.Context) (uuid.UUID, bool) {
	if ctx == nil {
		return uuid.Nil, false
	}

	value, ok := ctx.Get(ContextActorMemberIDKey)
	if !ok {
		return GetMemberID(ctx)
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

func GetActorIsAdmin(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}

	value, ok := ctx.Get(ContextActorIsAdminKey)
	if !ok {
		return GetIsAdmin(ctx)
	}

	isAdmin, ok := value.(bool)
	if !ok {
		return false
	}

	return isAdmin
}

func IsActingAs(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}

	value, ok := ctx.Get(ContextActingAsKey)
	if !ok {
		return false
	}

	isActingAs, ok := value.(bool)
	if !ok {
		return false
	}

	return isActingAs
}

func WithRequestMeta(ctx context.Context, endpoint string, actorMemberID uuid.UUID, targetMemberID uuid.UUID, isActingAs bool) context.Context {
	ctx = context.WithValue(ctx, requestEndpointKey, endpoint)
	ctx = context.WithValue(ctx, requestActorMemberKey, actorMemberID)
	ctx = context.WithValue(ctx, requestTargetMemberKey, targetMemberID)
	ctx = context.WithValue(ctx, requestActingAsKey, isActingAs)
	return ctx
}

func RequestEndpoint(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value := ctx.Value(requestEndpointKey)
	endpoint, _ := value.(string)
	return endpoint
}

func RequestActorMemberID(ctx context.Context) (uuid.UUID, bool) {
	if ctx == nil {
		return uuid.Nil, false
	}
	value := ctx.Value(requestActorMemberKey)
	memberID, ok := value.(uuid.UUID)
	return memberID, ok
}

func RequestTargetMemberID(ctx context.Context) (uuid.UUID, bool) {
	if ctx == nil {
		return uuid.Nil, false
	}
	value := ctx.Value(requestTargetMemberKey)
	memberID, ok := value.(uuid.UUID)
	return memberID, ok
}

func RequestIsActingAs(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	value := ctx.Value(requestActingAsKey)
	isActingAs, _ := value.(bool)
	return isActingAs
}
