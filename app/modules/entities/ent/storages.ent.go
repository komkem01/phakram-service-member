package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RelateTypeEnum string

const (
	RelateTypeMemberFile  RelateTypeEnum = "MEMBER_FILE"
	RelateTypeOrderFile   RelateTypeEnum = "ORDER_FILE"
	RelateTypeProductFile RelateTypeEnum = "PRODUCT_FILE"
	RelateTypePaymentFile RelateTypeEnum = "PAYMENT_FILE"
)

type StorageEntity struct {
	bun.BaseModel `bun:"table:storages"`

	ID            uuid.UUID      `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	RefID         uuid.UUID      `bun:"ref_id,type:uuid" json:"ref_id"`
	FileName      string         `bun:"file_name" json:"file_name"`
	FilePath      string         `bun:"file_path" json:"file_path"`
	FileType      string         `bun:"file_type" json:"file_type"`
	FileSize      string         `bun:"file_size" json:"file_size"`
	RelatedEntity RelateTypeEnum `bun:"related_entity" json:"related_entity"`
	UploadedBy    uuid.UUID      `bun:"uploaded_by,type:uuid" json:"uploaded_by"`
	CreatedAt     time.Time      `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
