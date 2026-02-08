package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ProductEntity struct {
	bun.BaseModel `bun:"table:products"`

	ID         uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	CategoryID uuid.UUID       `bun:"category_id,type:uuid" json:"category_id"`
	NameTh     string          `bun:"name_th" json:"name_th"`
	NameEn     string          `bun:"name_en" json:"name_en"`
	ProductNo  string          `bun:"product_no" json:"product_no"`
	Price      decimal.Decimal `bun:"price" json:"price"`
	IsActive   bool            `bun:"is_active" json:"is_active"`
	CreatedAt  time.Time       `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time       `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	DeletedAt  *time.Time      `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at"`
}
