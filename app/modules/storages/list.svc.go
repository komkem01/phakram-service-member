package storages

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListStorageServiceRequest struct {
	base.RequestPaginate
}

type ListStorageServiceResponses struct {
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

func (s *Service) ListService(ctx context.Context, req *ListStorageServiceRequest) ([]*ListStorageServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.list.start`)

	data, page, err := s.db.ListStorages(ctx, &entitiesdto.ListStoragesRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListStorageServiceResponses
	for _, item := range data {
		temp := &ListStorageServiceResponses{
			ID:            item.ID,
			RefID:         item.RefID,
			FileName:      item.FileName,
			FilePath:      item.FilePath,
			FileType:      item.FileType,
			FileSize:      item.FileSize,
			RelatedEntity: string(item.RelatedEntity),
			UploadedBy:    item.UploadedBy,
			CreatedAt:     item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`storages.svc.list.copy`)
	return response, page, nil
}
