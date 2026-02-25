package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ZipcodeEntity struct {
	bun.BaseModel `bun:"table:zipcodes"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	SubDistrictsID uuid.UUID `bun:"sub_districts_id,type:uuid" json:"sub_districts_id"`
	Name           string    `bun:"name" json:"name"`
	IsActive       bool      `bun:"is_active" json:"is_active"`
	CreatedAt      time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
