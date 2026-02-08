package provinces

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateProvinceService struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateProvinceService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`provinces.svc.update.start`)

	data, err := s.db.GetProvinceByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.Name != "" {
		data.Name = req.Name
	}
	if req.IsActive != nil {
		data.IsActive = *req.IsActive
	}

	if err := s.db.UpdateProvince(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`provinces.svc.update.success`)
	return nil
}
