package statuses

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoStatusServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoStatusServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`statuses.svc.info.start`)

	data, err := s.db.GetStatusByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoStatusServiceResponses{
		ID:        data.ID,
		NameTh:    data.NameTh,
		NameEn:    data.NameEn,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`statuses.svc.info.success`)
	return resp, nil
}
