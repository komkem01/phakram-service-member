package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrderCancellationEntity struct {
	bun.BaseModel `bun:"table:order_cancellations"`

	ID            uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	OrderID       uuid.UUID  `bun:"order_id,type:uuid" json:"order_id"`
	CancelledBy   *uuid.UUID `bun:"cancelled_by,type:uuid" json:"cancelled_by"`
	CancelledRole string     `bun:"cancelled_role" json:"cancelled_role"`
	Reason        string     `bun:"reason" json:"reason"`
	CreatedAt     time.Time  `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
