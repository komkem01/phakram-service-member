package storages

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateStorageService struct {
	RefID         uuid.UUID `json:"ref_id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileType      string    `json:"file_type"`
	FileSize      string    `json:"file_size"`
	RelatedEntity string    `json:"related_entity"`
	UploadedBy    uuid.UUID `json:"uploaded_by"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateStorageService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.update.start`)

	data, err := s.db.GetStorageByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.RefID != uuid.Nil {
		data.RefID = req.RefID
	}
	if req.FileName != "" {
		data.FileName = req.FileName
	}
	if req.FilePath != "" {
		data.FilePath = req.FilePath
	}
	if req.FileType != "" {
		data.FileType = req.FileType
	}
	if req.FileSize != "" {
		data.FileSize = req.FileSize
	}
	if req.RelatedEntity != "" {
		data.RelatedEntity = ent.RelateTypeEnum(req.RelatedEntity)
	}
	if req.UploadedBy != uuid.Nil {
		data.UploadedBy = req.UploadedBy
	}

	if err := s.db.UpdateStorage(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`storages.svc.update.success`)
	return nil
}
