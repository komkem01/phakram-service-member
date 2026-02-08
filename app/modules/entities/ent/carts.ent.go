package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CartEntity struct {
	bun.BaseModel `bun:"table:carts"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID  uuid.UUID `bun:"member_id,type:uuid" json:"member_id"`
	IsActive  bool      `bun:"is_active" json:"is_active"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
