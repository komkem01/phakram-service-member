package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ProductStockEntity struct {
	bun.BaseModel `bun:"table:product_stocks"`

	ID          uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ProductID   uuid.UUID       `bun:"product_id,type:uuid" json:"product_id"`
	UnitPrice   decimal.Decimal `bun:"unit_price" json:"unit_price"`
	StockAmount int             `bun:"stock_amount" json:"stock_amount"`
	Remaining   int             `bun:"remaining" json:"remaining"`
	CreatedAt   time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	DeletedAt   *time.Time      `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at"`
}
