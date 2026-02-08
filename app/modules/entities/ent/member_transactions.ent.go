package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ActionTypeEnum represents member transaction actions.
type ActionTypeEnum string

const (
	ActionTypeCreated    ActionTypeEnum = "created"
	ActionTypeUpdated    ActionTypeEnum = "updated"
	ActionTypeDeleted    ActionTypeEnum = "deleted"
	ActionTypeLogined    ActionTypeEnum = "logined"
	ActionTypeRegistered ActionTypeEnum = "registered"
)

type MemberTransactionEntity struct {
	bun.BaseModel `bun:"table:member_transactions"`

	ID        uuid.UUID      `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID  uuid.UUID      `bun:"member_id,type:uuid" json:"member_id"`
	Action    ActionTypeEnum `bun:"action" json:"action"`
	Details   string         `bun:"details" json:"details"`
	CreatedAt time.Time      `bun:"created_at,default:current_timestamp" json:"created_at"`
}
