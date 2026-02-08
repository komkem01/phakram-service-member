package storages

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateStorageService struct {
	RefID         uuid.UUID `json:"ref_id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileType      string    `json:"file_type"`
	FileSize      string    `json:"file_size"`
	RelatedEntity string    `json:"related_entity"`
	UploadedBy    uuid.UUID `json:"uploaded_by"`
}

func (s *Service) CreateStorageService(ctx context.Context, req *CreateStorageService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.create.start`)

	storage := &ent.StorageEntity{
		ID:            uuid.New(),
		RefID:         req.RefID,
		FileName:      req.FileName,
		FilePath:      req.FilePath,
		FileType:      req.FileType,
		FileSize:      req.FileSize,
		RelatedEntity: ent.RelateTypeEnum(req.RelatedEntity),
		UploadedBy:    req.UploadedBy,
	}
	if err := s.db.CreateStorage(ctx, storage); err != nil {
		return err
	}
	span.AddEvent(`storages.svc.create.success`)
	return nil
}
