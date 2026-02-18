package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrderPaymentReviewEntity struct {
	bun.BaseModel `bun:"table:order_payment_reviews"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	OrderID        uuid.UUID  `bun:"order_id,type:uuid" json:"order_id"`
	PaymentID      uuid.UUID  `bun:"payment_id,type:uuid" json:"payment_id"`
	ReviewStatus   string     `bun:"review_status" json:"review_status"`
	RejectedReason string     `bun:"rejected_reason" json:"rejected_reason"`
	ReviewedBy     *uuid.UUID `bun:"reviewed_by,type:uuid" json:"reviewed_by"`
	ReviewedAt     *time.Time `bun:"reviewed_at" json:"reviewed_at"`
	CreatedAt      time.Time  `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
