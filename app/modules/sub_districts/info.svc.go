package sub_districts

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoSubDistrictServiceResponses struct {
	ID         uuid.UUID `json:"id"`
	DistrictID uuid.UUID `json:"district_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoSubDistrictServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`sub_districts.svc.info.start`)

	data, err := s.db.GetSubDistrictByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoSubDistrictServiceResponses{
		ID:         data.ID,
		DistrictID: data.DistrictID,
		Name:       data.Name,
		IsActive:   data.IsActive,
		CreatedAt:  data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`sub_districts.svc.info.success`)
	return resp, nil
}
