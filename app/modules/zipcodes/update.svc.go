package zipcodes

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateZipcodeService struct {
	SubDistrictsID *uuid.UUID `json:"sub_districts_id"`
	Name           string     `json:"name"`
	IsActive       *bool      `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateZipcodeService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`zipcodes.svc.update.start`)

	data, err := s.db.GetZipcodeByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.SubDistrictsID != nil {
		data.SubDistrictsID = *req.SubDistrictsID
	}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.IsActive != nil {
		data.IsActive = *req.IsActive
	}

	if err := s.db.UpdateZipcode(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`zipcodes.svc.update.success`)
	return nil
}
