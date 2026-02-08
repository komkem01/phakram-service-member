package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StatusEntity struct {
	bun.BaseModel `bun:"table:statuses"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	NameTh    string    `bun:"name_th" json:"name_th"`
	NameEn    string    `bun:"name_en" json:"name_en"`
	IsActive  bool      `bun:"is_active,default:false" json:"is_active"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
}
