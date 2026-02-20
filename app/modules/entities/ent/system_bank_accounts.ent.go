package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SystemBankAccountEntity struct {
	bun.BaseModel `bun:"table:system_bank_accounts"`

	ID               uuid.UUID `bun:"id,pk,type:uuid" json:"id"`
	BankID           uuid.UUID `bun:"bank_id,type:uuid" json:"bank_id"`
	AccountName      string    `bun:"account_name" json:"account_name"`
	AccountNo        string    `bun:"account_no" json:"account_no"`
	Branch           string    `bun:"branch" json:"branch"`
	QRCodeImageURL   string    `bun:"qr_image_url" json:"qr_image_url"`
	IsActive         bool      `bun:"is_active" json:"is_active"`
	IsDefaultReceive bool      `bun:"is_default_receive" json:"is_default_receive"`
	IsDefaultRefund  bool      `bun:"is_default_refund" json:"is_default_refund"`
	CreatedAt        time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt        time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
