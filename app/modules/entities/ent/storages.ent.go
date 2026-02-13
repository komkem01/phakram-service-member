package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RelatedEntityEnum string

const (
	RelatedEntityMemberFile  RelatedEntityEnum = "MEMBER_FILE"
	RelatedEntityOrderFile   RelatedEntityEnum = "ORDER_FILE"
	RelatedEntityProductFile RelatedEntityEnum = "PRODUCT_FILE"
	RelatedEntityPaymentFile RelatedEntityEnum = "PAYMENT_FILE"
	RelatedEntityOther       RelatedEntityEnum = "OTHER"
)

type StorageEntity struct {
	bun.BaseModel `bun:"table:storages"`

	ID            uuid.UUID         `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	RefID         uuid.UUID         `bun:"ref_id,type:uuid" json:"ref_id"`
	FileName      string            `bun:"file_name" json:"file_name"`
	FilePath      string            `bun:"file_path" json:"file_path"`
	FileSize      int64             `bun:"file_size" json:"file_size"`
	FileType      string            `bun:"file_type" json:"file_type"`
	IsActive      bool              `bun:"is_active" json:"is_active"`
	RelatedEntity RelatedEntityEnum `bun:"related_entity" json:"related_entity"`
	UploadedBy    uuid.UUID         `bun:"uploaded_by,type:uuid" json:"uploaded_by"`
	CreatedAt     time.Time         `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time         `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
