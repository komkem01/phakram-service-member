package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberBankEntity struct {
	bun.BaseModel `bun:"table:member_banks"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID    uuid.UUID `bun:"member_id,type:uuid" json:"member_id"`
	BankID      uuid.UUID `bun:"bank_id,type:uuid" json:"bank_id"`
	BankNo      string    `bun:"bank_no" json:"bank_no"`
	FirstnameTh string    `bun:"firstname_th" json:"firstname_th"`
	LastnameTh  string    `bun:"lastname_th" json:"lastname_th"`
	FirstnameEn string    `bun:"firstname_en" json:"firstname_en"`
	LastnameEn  string    `bun:"lastname_en" json:"lastname_en"`
	IsDefault   bool      `bun:"is_default" json:"is_default"`
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
