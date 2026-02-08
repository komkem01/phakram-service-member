package product_files

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListProductFileServiceRequest struct {
	base.RequestPaginate
}

type ListProductFileServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListProductFileServiceRequest) ([]*ListProductFileServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_files.svc.list.start`)

	data, page, err := s.db.ListProductFiles(ctx, &entitiesdto.ListProductFilesRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListProductFileServiceResponses
	for _, item := range data {
		temp := &ListProductFileServiceResponses{
			ID:        item.ID,
			ProductID: item.ProductID,
			FileID:    item.FileID,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`product_files.svc.list.copy`)
	return response, page, nil
}
