package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"phakram/app/modules"
	"phakram/app/modules/auth"
	"phakram/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func auditFailureMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		status := ctx.Writer.Status()
		fullPath := ctx.FullPath()
		if fullPath == "" || strings.HasPrefix(fullPath, "/healthz") {
			return
		}
		if !strings.HasPrefix(fullPath, "/api/") {
			return
		}

		action := auditActionForRequest(ctx)
		actionType := strings.ToLower(ctx.Request.Method) + " " + fullPath
		actionID := auditActionID(ctx)
		actionBy := auditActionBy(ctx)

		auditStatus := ent.StatusAuditSuccesses
		detail := fmt.Sprintf("Request success: %s %s (%d)", ctx.Request.Method, fullPath, status)
		if status >= http.StatusBadRequest {
			auditStatus = ent.StatusAuditFailed
			detail = fmt.Sprintf("Request failed: %s %s (%d)", ctx.Request.Method, fullPath, status)
		}

		log := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       action,
			ActionType:   actionType,
			ActionID:     actionID,
			ActionBy:     actionBy,
			Status:       auditStatus,
			ActionDetail: detail,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_ = mod.ENT.Svc.CreateAuditLog(ctx.Request.Context(), log)

		memberAction, ok := memberActionForAuditAction(action)
		if !ok {
			return
		}

		actorMemberID, ok := auth.GetActorMemberID(ctx)
		if !ok || actorMemberID == uuid.Nil {
			return
		}

		memberTransaction := &ent.MemberTransactionEntity{
			ID:       uuid.New(),
			MemberID: actorMemberID,
			Action:   memberAction,
			Details:  detail,
		}
		_ = mod.ENT.Svc.CreateMemberTransaction(ctx.Request.Context(), memberTransaction)
	}
}

func memberActionForAuditAction(action ent.AuditActionEnum) (ent.MemberActionEnum, bool) {
	switch action {
	case ent.AuditActionCreated:
		return ent.MemberActionCreated, true
	case ent.AuditActionUpdated:
		return ent.MemberActionUpdated, true
	case ent.AuditActionDeleted:
		return ent.MemberActionDeleted, true
	case ent.AuditActionLogined:
		return ent.MemberActionLogined, true
	case ent.AuditActionRegistered:
		return ent.MemberActionRegistered, true
	case ent.AuditActionRead:
		return ent.MemberActionRead, true
	default:
		return "", false
	}
}

func auditActionForRequest(ctx *gin.Context) ent.AuditActionEnum {
	path := strings.ToLower(ctx.FullPath())
	method := strings.ToUpper(ctx.Request.Method)

	if strings.Contains(path, "/login") {
		return ent.AuditActionLogined
	}
	if strings.Contains(path, "/register") {
		return ent.AuditActionRegistered
	}

	switch method {
	case http.MethodPost:
		return ent.AuditActionCreated
	case http.MethodDelete:
		return ent.AuditActionDeleted
	case http.MethodPut, http.MethodPatch:
		return ent.AuditActionUpdated
	case http.MethodGet:
		return ent.AuditActionRead
	default:
		return ent.AuditActionUpdated
	}
}

func auditActionID(ctx *gin.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	if idParam := ctx.Param("id"); idParam != "" {
		if id, err := uuid.Parse(idParam); err == nil {
			return id
		}
	}

	return uuid.Nil
}

func auditActionBy(ctx *gin.Context) *uuid.UUID {
	if ctx == nil {
		return nil
	}

	memberID, ok := auth.GetActorMemberID(ctx)
	if !ok || memberID == uuid.Nil {
		return nil
	}

	return &memberID
}
