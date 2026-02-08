package sub_districts

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateSubDistrictService struct {
	DistrictID *uuid.UUID `json:"district_id"`
	Name       string     `json:"name"`
	IsActive   *bool      `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateSubDistrictService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`sub_districts.svc.update.start`)

	data, err := s.db.GetSubDistrictByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.DistrictID != nil {
		data.DistrictID = *req.DistrictID
	}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.IsActive != nil {
		data.IsActive = *req.IsActive
	}

	if err := s.db.UpdateSubDistrict(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`sub_districts.svc.update.success`)
	return nil
}
