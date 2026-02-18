package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type StatusTypeEnum string

const (
	StatusTypePending         StatusTypeEnum = "pending"
	StatusTypePaid            StatusTypeEnum = "paid"
	StatusTypeRefundRequested StatusTypeEnum = "refund_requested"
	StatusTypeShipping        StatusTypeEnum = "shipping"
	StatusTypeCompleted       StatusTypeEnum = "completed"
	StatusTypeCancelled       StatusTypeEnum = "cancelled"
)

type OrderEntity struct {
	bun.BaseModel `bun:"table:orders"`

	ID                     uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	OrderNo                string          `bun:"order_no" json:"order_no"`
	MemberID               uuid.UUID       `bun:"member_id,type:uuid" json:"member_id"`
	PaymentID              uuid.UUID       `bun:"payment_id,type:uuid" json:"payment_id"`
	AddressID              uuid.UUID       `bun:"address_id,type:uuid" json:"address_id"`
	Status                 StatusTypeEnum  `bun:"status" json:"status"`
	TotalAmount            decimal.Decimal `bun:"total_amount" json:"total_amount"`
	DiscountAmount         decimal.Decimal `bun:"discount_amount" json:"discount_amount"`
	NetAmount              decimal.Decimal `bun:"net_amount" json:"net_amount"`
	CreatedAt              time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt              time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	PaymentSubmitted       bool            `bun:"-" json:"payment_submitted"`
	PaymentRejected        bool            `bun:"-" json:"payment_rejected"`
	PaymentRejectionReason string          `bun:"-" json:"payment_rejection_reason,omitempty"`
	PaymentAppealReason    string          `bun:"-" json:"payment_appeal_reason,omitempty"`
	RefundRejectionReason  string          `bun:"-" json:"refund_rejection_reason,omitempty"`
	CancellationReason     string          `bun:"-" json:"cancellation_reason,omitempty"`
	ShippingTrackingNo     string          `bun:"-" json:"shipping_tracking_no,omitempty"`
	StatusSummary          string          `bun:"-" json:"status_summary,omitempty"`
	StatusNextStep         string          `bun:"-" json:"status_next_step,omitempty"`
}
