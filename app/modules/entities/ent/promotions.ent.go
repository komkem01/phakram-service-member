package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type promotionRecord struct {
	bun.BaseModel `bun:"table:promotions"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid"`
	Code           string     `bun:"code,notnull"`
	Name           string     `bun:"name,notnull"`
	Description    string     `bun:"description"`
	DiscountType   string     `bun:"discount_type,notnull"`
	DiscountValue  float64    `bun:"discount_value,notnull"`
	MaxDiscount    *float64   `bun:"max_discount"`
	MinOrderAmount float64    `bun:"min_order_amount,notnull"`
	UsageLimit     *int       `bun:"usage_limit"`
	UsagePerMember *int       `bun:"usage_per_member"`
	UsedCount      int        `bun:"used_count,notnull"`
	StartsAt       *time.Time `bun:"starts_at"`
	EndsAt         *time.Time `bun:"ends_at"`
	IsActive       bool       `bun:"is_active,notnull"`
	CreatedAt      time.Time  `bun:"created_at,notnull"`
	UpdatedAt      time.Time  `bun:"updated_at,notnull"`
}
