package zipcodes

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoZipcodeServiceResponses struct {
	ID             uuid.UUID `json:"id"`
	SubDistrictsID uuid.UUID `json:"sub_districts_id"`
	Name           string    `json:"name"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoZipcodeServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`zipcodes.svc.info.start`)

	data, err := s.db.GetZipcodeByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoZipcodeServiceResponses{
		ID:             data.ID,
		SubDistrictsID: data.SubDistrictsID,
		Name:           data.Name,
		IsActive:       data.IsActive,
		CreatedAt:      data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`zipcodes.svc.info.success`)
	return resp, nil
}
