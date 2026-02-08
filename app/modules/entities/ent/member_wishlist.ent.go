package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type MemberWishlistEntity struct {
	bun.BaseModel `bun:"table:member_wishlist"`

	ID              uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID        uuid.UUID       `bun:"member_id,type:uuid" json:"member_id"`
	ProductID       uuid.UUID       `bun:"product_id,type:uuid" json:"product_id"`
	Quantity        int             `bun:"quantity" json:"quantity"`
	PricePerUnit    decimal.Decimal `bun:"price_per_unit" json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `bun:"total_item_amount" json:"total_item_amount"`
	CreatedAt       time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
