package storages

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoStorageServiceResponses struct {
	ID            uuid.UUID `json:"id"`
	RefID         uuid.UUID `json:"ref_id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileType      string    `json:"file_type"`
	FileSize      string    `json:"file_size"`
	RelatedEntity string    `json:"related_entity"`
	UploadedBy    uuid.UUID `json:"uploaded_by"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoStorageServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.info.start`)

	data, err := s.db.GetStorageByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoStorageServiceResponses{
		ID:            data.ID,
		RefID:         data.RefID,
		FileName:      data.FileName,
		FilePath:      data.FilePath,
		FileType:      data.FileType,
		FileSize:      data.FileSize,
		RelatedEntity: string(data.RelatedEntity),
		UploadedBy:    data.UploadedBy,
		CreatedAt:     data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`storages.svc.info.success`)
	return resp, nil
}
