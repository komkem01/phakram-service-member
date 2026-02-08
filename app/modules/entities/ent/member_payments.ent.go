package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type MemberPaymentEntity struct {
	bun.BaseModel `bun:"table:member_payments"`

	ID        uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID  uuid.UUID       `bun:"member_id,type:uuid" json:"member_id"`
	PaymentID uuid.UUID       `bun:"payment_id,type:uuid" json:"payment_id"`
	Quantity  int             `bun:"quantity" json:"quantity"`
	Price     decimal.Decimal `bun:"price" json:"price"`
	CreatedAt time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
