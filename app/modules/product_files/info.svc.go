package product_files

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoProductFileServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoProductFileServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_files.svc.info.start`)

	data, err := s.db.GetProductFileByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoProductFileServiceResponses{
		ID:        data.ID,
		ProductID: data.ProductID,
		FileID:    data.FileID,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`product_files.svc.info.success`)
	return resp, nil
}
