package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberAccountEntity struct {
	bun.BaseModel `bun:"table:member_accounts"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID  uuid.UUID `bun:"member_id,type:uuid" json:"member_id"`
	Email     string    `bun:"email" json:"email"`
	Password  string    `bun:"password" json:"password"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
