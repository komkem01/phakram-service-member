package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductFileEntity struct {
	bun.BaseModel `bun:"table:product_files"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ProductID uuid.UUID  `bun:"product_id,type:uuid" json:"product_id"`
	FileID    uuid.UUID  `bun:"file_id,type:uuid" json:"file_id"`
	CreatedAt time.Time  `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at"`
}
