package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type RoleTypeEnum string

const (
	RoleTypeCustomer RoleTypeEnum = "customer"
	RoleTypeAdmin    RoleTypeEnum = "admin"
)

type MemberEntity struct {
	bun.BaseModel `bun:"table:members"`

	ID            uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberNo      string          `bun:"member_no" json:"member_no"`
	TierID        uuid.UUID       `bun:"tier_id,type:uuid" json:"tier_id"`
	StatusID      uuid.UUID       `bun:"status_id,type:uuid" json:"status_id"`
	PrefixID      uuid.UUID       `bun:"prefix_id,type:uuid" json:"prefix_id"`
	GenderID      uuid.UUID       `bun:"gender_id,type:uuid" json:"gender_id"`
	FirstnameTh   string          `bun:"firstname_th" json:"firstname_th"`
	LastnameTh    string          `bun:"lastname_th" json:"lastname_th"`
	FirstnameEn   string          `bun:"firstname_en" json:"firstname_en"`
	LastnameEn    string          `bun:"lastname_en" json:"lastname_en"`
	Role          RoleTypeEnum    `bun:"role" json:"role"`
	Phone         string          `bun:"phone" json:"phone"`
	TotalSpent    decimal.Decimal `bun:"total_spent" json:"total_spent"`
	CurrentPoints int             `bun:"current_points" json:"current_points"`
	Registration  *time.Time      `bun:"registration" json:"registration"`
	LastLogin     *time.Time      `bun:"last_login" json:"last_login"`
	CreatedAt     time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	DeletedAt     *time.Time      `bun:"deleted_at,soft_delete" json:"deleted_at"`
}
