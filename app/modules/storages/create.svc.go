package storages

import (
	"context"
	"strings"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateStorageServiceRequest struct {
	RefID         uuid.UUID
	FileName      string
	FilePath      string
	FileSize      int64
	FileType      string
	RelatedEntity string
	UploadedBy    uuid.UUID
	IsActive      *bool
}

func (s *Service) CreateService(ctx context.Context, req *CreateStorageServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.create.start`)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	storage := &ent.StorageEntity{
		ID:            uuid.New(),
		RefID:         req.RefID,
		FileName:      req.FileName,
		FilePath:      req.FilePath,
		FileSize:      req.FileSize,
		FileType:      req.FileType,
		RelatedEntity: parseRelatedEntity(req.RelatedEntity),
		UploadedBy:    req.UploadedBy,
		IsActive:      isActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.UploadStorage(ctx, storage); err != nil {
		return err
	}

	span.AddEvent(`storages.svc.create.success`)
	return nil
}

func parseRelatedEntity(input string) ent.RelatedEntityEnum {
	value := strings.ToUpper(strings.TrimSpace(input))
	switch value {
	case string(ent.RelatedEntityMemberFile):
		return ent.RelatedEntityMemberFile
	case string(ent.RelatedEntityOrderFile):
		return ent.RelatedEntityOrderFile
	case string(ent.RelatedEntityProductFile):
		return ent.RelatedEntityProductFile
	case string(ent.RelatedEntityPaymentFile):
		return ent.RelatedEntityPaymentFile
	default:
		return ent.RelatedEntityOther
	}
}
