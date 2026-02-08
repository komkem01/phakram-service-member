package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CategoryEntity struct {
	bun.BaseModel `bun:"table:categories"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ParentID  *uuid.UUID `bun:"parent_id,type:uuid" json:"parent_id"`
	NameTh    string     `bun:"name_th" json:"name_th"`
	NameEn    string     `bun:"name_en" json:"name_en"`
	IsActive  bool       `bun:"is_active" json:"is_active"`
	CreatedAt time.Time  `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
