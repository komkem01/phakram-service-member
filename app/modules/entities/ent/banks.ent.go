package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BankEntity struct {
	bun.BaseModel `bun:"table:banks"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	NameTh    string    `bun:"name_th" json:"name_th"`
	NameAbbTh string    `bun:"name_abb_th" json:"name_abb_th"`
	NameEn    string    `bun:"name_en" json:"name_en"`
	NameAbbEn string    `bun:"name_abb_en" json:"name_abb_en"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
