package districts

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoDistrictServiceResponses struct {
	ID         uuid.UUID `json:"id"`
	ProvinceID uuid.UUID `json:"province_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoDistrictServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`districts.svc.info.start`)

	data, err := s.db.GetDistrictByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoDistrictServiceResponses{
		ID:         data.ID,
		ProvinceID: data.ProvinceID,
		Name:       data.Name,
		IsActive:   data.IsActive,
		CreatedAt:  data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`districts.svc.info.success`)
	return resp, nil
}
