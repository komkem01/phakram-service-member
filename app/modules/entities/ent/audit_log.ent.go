package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StatusAuditEnum string

type AuditActionEnum string

const (
	StatusAuditSuccesses StatusAuditEnum = "successes"
	StatusAuditFailed    StatusAuditEnum = "failed"
)

const (
	AuditActionCreated    AuditActionEnum = "created"
	AuditActionUpdated    AuditActionEnum = "updated"
	AuditActionDeleted    AuditActionEnum = "deleted"
	AuditActionLogined    AuditActionEnum = "logined"
	AuditActionRegistered AuditActionEnum = "registered"
)

type AuditLogEntity struct {
	bun.BaseModel `bun:"table:audit_log"`

	ID           uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Action       AuditActionEnum `bun:"action" json:"action"`
	ActionType   string          `bun:"action_type" json:"action_type"`
	ActionID     uuid.UUID       `bun:"action_id,type:uuid" json:"action_id"`
	ActionBy     uuid.UUID       `bun:"action_by,type:uuid" json:"action_by"`
	Status       StatusAuditEnum `bun:"status" json:"status"`
	ActionDetail string          `bun:"action_detail" json:"action_detail"`
	CreatedAt    time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
