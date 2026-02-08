package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ActionAuditEnum represents audit actions.
type ActionAuditEnum string

const (
	ActionAuditCreate  ActionAuditEnum = "create"
	ActionAuditUpdate  ActionAuditEnum = "update"
	ActionAuditDelete  ActionAuditEnum = "delete"
	ActionAuditOrder   ActionAuditEnum = "order"
	ActionAuditPayment ActionAuditEnum = "payment"
)

// StatusAuditEnum represents audit status.
type StatusAuditEnum string

const (
	StatusAuditSuccess StatusAuditEnum = "success"
	StatusAuditFail    StatusAuditEnum = "fail"
)

type AuditLogEntity struct {
	bun.BaseModel `bun:"table:audit_log"`

	ID           uuid.UUID       `bun:"id,pk,type:uuid" json:"id"`
	Action       ActionAuditEnum `bun:"action" json:"action"`
	ActionType   string          `bun:"action_type" json:"action_type"`
	ActionID     *uuid.UUID      `bun:"action_id,type:uuid" json:"action_id"`
	ActionBy     *uuid.UUID      `bun:"action_by,type:uuid" json:"action_by"`
	Status       StatusAuditEnum `bun:"status" json:"status"`
	ActionDetail string          `bun:"action_detail" json:"action_detail"`
	CreatedAt    time.Time       `bun:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `bun:"updated_at" json:"updated_at"`
}
