package ent

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ProductDetailEntity struct {
	bun.BaseModel `bun:"table:product_details"`

	ID               uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ProductID        uuid.UUID       `bun:"product_id,type:uuid" json:"product_id"`
	Description      string          `bun:"description" json:"description"`
	Material         string          `bun:"material" json:"material"`
	Dimensions       string          `bun:"dimensions" json:"dimensions"`
	Weight           decimal.Decimal `bun:"weight" json:"weight"`
	CareInstructions string          `bun:"care_instructions" json:"care_instructions"`
}
