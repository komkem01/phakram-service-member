package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DistrictEntity struct {
	bun.BaseModel `bun:"table:districts"`

	ID         uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ProvinceID uuid.UUID `bun:"province_id,type:uuid" json:"province_id"`
	Name       string    `bun:"name" json:"name"`
	IsActive   bool      `bun:"is_active" json:"is_active"`
	CreatedAt  time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
