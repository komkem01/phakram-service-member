package statuses

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateStatusService struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateStatusService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`statuses.svc.update.start`)

	data, err := s.db.GetStatusByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.NameTh != "" {
		data.NameTh = req.NameTh
	}
	if req.NameEn != "" {
		data.NameEn = req.NameEn
	}

	if err := s.db.UpdateStatus(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`statuses.svc.update.success`)
	return nil
}
