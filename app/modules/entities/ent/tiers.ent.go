package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type TierEntity struct {
	bun.BaseModel `bun:"table:tiers"`

	ID           uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	NameTh       string          `bun:"name_th" json:"name_th"`
	NameEn       string          `bun:"name_en" json:"name_en"`
	MinSpending  decimal.Decimal `bun:"min_spending" json:"min_spending"`
	IsActive     bool            `bun:"is_active" json:"is_active"`
	DiscountRate decimal.Decimal `bun:"discount_rate" json:"discount_rate"`
	CreatedAt    time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
