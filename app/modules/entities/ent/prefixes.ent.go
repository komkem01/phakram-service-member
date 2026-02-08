package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PrefixEntity struct {
	bun.BaseModel `bun:"table:prefixes"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	NameTh    string    `bun:"name_th" json:"name_th"`
	NameEn    string    `bun:"name_en" json:"name_en"`
	GenderID  uuid.UUID `bun:"gender_id,type:uuid" json:"gender_id"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
}
