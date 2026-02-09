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
		if status < http.StatusBadRequest {
			return
		}

		fullPath := ctx.FullPath()
		if fullPath == "" || strings.HasPrefix(fullPath, "/healthz") {
			return
		}

		action := auditActionForRequest(ctx)
		actionType := strings.ToLower(ctx.Request.Method) + " " + fullPath
		actionID := auditActionID(ctx)
		actionBy := auditActionBy(ctx)

		log := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       action,
			ActionType:   actionType,
			ActionID:     actionID,
			ActionBy:     actionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Request failed: %s %s (%d)", ctx.Request.Method, fullPath, status),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_ = mod.ENT.Svc.CreateAuditLog(ctx.Request.Context(), log)
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

	memberID, ok := auth.GetMemberID(ctx)
	if !ok || memberID == uuid.Nil {
		return nil
	}

	return &memberID
}
