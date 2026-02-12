package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberFileEntity struct {
	bun.BaseModel `bun:"table:member_files"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID  uuid.UUID  `bun:"member_id,type:uuid" json:"member_id"`
	FileID    uuid.UUID  `bun:"file_id,type:uuid" json:"file_id"`
	CreatedAt time.Time  `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete" json:"deleted_at"`
}
