package genders

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateGenderService struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
	IsActive bool   `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateGenderService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`genders.svc.update.start`)

	data, err := s.db.GetGenderByID(ctx, id)
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

	data.IsActive = req.IsActive

	if err := s.db.UpdateGender(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`genders.svc.update.success`)
	return nil
}
