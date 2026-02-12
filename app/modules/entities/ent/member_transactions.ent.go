package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberActionEnum string

const (
	MemberActionCreated    MemberActionEnum = "created"
	MemberActionUpdated    MemberActionEnum = "updated"
	MemberActionDeleted    MemberActionEnum = "deleted"
	MemberActionLogined    MemberActionEnum = "logined"
	MemberActionRegistered MemberActionEnum = "registered"
)

type MemberTransactionEntity struct {
	bun.BaseModel `bun:"table:member_transactions"`

	ID        uuid.UUID        `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID  uuid.UUID        `bun:"member_id,type:uuid" json:"member_id"`
	Action    MemberActionEnum `bun:"action" json:"action"`
	Details   string           `bun:"details" json:"details"`
	CreatedAt time.Time        `bun:"created_at,default:current_timestamp" json:"created_at"`
}
