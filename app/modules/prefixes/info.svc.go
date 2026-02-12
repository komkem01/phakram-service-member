package prefixes

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"phakram/app/utils"
	"phakram/config/i18n"

	"github.com/google/uuid"
)

type InfoPrefixServiceResponses struct {
	ID         uuid.UUID `json:"id"`
	NameTh     string    `json:"name_th"`
	NameEn     string    `json:"name_en"`
	GenderID   uuid.UUID `json:"gender_id"`
	GenderName string    `json:"gender_name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoPrefixServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`prefixes.svc.info.start`)

	data, err := s.db.GetPrefixByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.With(slog.Any(`id`, id)).Errf(`prefix not found: %s`, err)
			return nil, i18n.ErrPrefixNotFound
		}
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	genderData, err := s.dbGender.GetGenderByID(ctx, data.GenderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.With(slog.Any(`gender_id`, data.GenderID)).Errf(`gender not found: %s`, err)
			return nil, i18n.ErrGenderNotFound
		}
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoPrefixServiceResponses{
		ID:         data.ID,
		NameTh:     data.NameTh,
		NameEn:     data.NameEn,
		GenderID:   data.GenderID,
		GenderName: genderData.NameTh,
		IsActive:   data.IsActive,
		CreatedAt:  data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`prefixes.svc.info.success`)
	return resp, nil
}
