package districts

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateDistrictService struct {
	ProvinceID *uuid.UUID `json:"province_id"`
	Name       string     `json:"name"`
	IsActive   *bool      `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateDistrictService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`districts.svc.update.start`)

	data, err := s.db.GetDistrictByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.ProvinceID != nil {
		data.ProvinceID = *req.ProvinceID
	}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.IsActive != nil {
		data.IsActive = *req.IsActive
	}

	if err := s.db.UpdateDistrict(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`districts.svc.update.success`)
	return nil
}
