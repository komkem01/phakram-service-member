package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// PaymentTypeEnum represents payment status types.
type PaymentTypeEnum string

const (
	PaymentTypePending  PaymentTypeEnum = "pending"
	PaymentTypeSuccess  PaymentTypeEnum = "success"
	PaymentTypeFailed   PaymentTypeEnum = "failed"
	PaymentTypeRefunded PaymentTypeEnum = "refunded"
)

type PaymentEntity struct {
	bun.BaseModel `bun:"table:payments"`

	ID         uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Amount     decimal.Decimal `bun:"amount" json:"amount"`
	Status     PaymentTypeEnum `bun:"status" json:"status"`
	ApprovedBy uuid.UUID       `bun:"approved_by,type:uuid" json:"approved_by"`
	ApprovedAt time.Time       `bun:"approved_at" json:"approved_at"`
}
