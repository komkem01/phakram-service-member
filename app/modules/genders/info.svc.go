package genders

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoGenderServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	CreatedAt string    `json:"created_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoGenderServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`genders.svc.info.start`)

	data, err := s.db.GetGenderByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoGenderServiceResponses{
		ID:        data.ID,
		NameTh:    data.NameTh,
		NameEn:    data.NameEn,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`genders.svc.info.success`)
	return resp, nil
}
