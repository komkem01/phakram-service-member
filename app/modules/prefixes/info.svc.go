package prefixes

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoPrefixServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	GenderID  uuid.UUID
	CreatedAt string    `json:"created_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoPrefixServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`prefixes.svc.info.start`)

	data, err := s.db.GetPrefixByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoPrefixServiceResponses{
		ID:        data.ID,
		NameTh:    data.NameTh,
		NameEn:    data.NameEn,
		GenderID:  data.GenderID,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`prefixes.svc.info.success`)
	return resp, nil
}
