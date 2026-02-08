package prefixes

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdatePrefixService struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
	GenderID uuid.UUID `json:"gender_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdatePrefixService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`prefixes.svc.update.start`)

	data, err := s.db.GetPrefixByID(ctx, id)
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
	if req.GenderID != uuid.Nil {
		data.GenderID = req.GenderID
	}

	if err := s.db.UpdatePrefix(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`prefixes.svc.update.success`)
	return nil
}
