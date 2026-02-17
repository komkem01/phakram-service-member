package banks

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoBankServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameAbbTh string    `json:"name_abb_th"`
	NameEn    string    `json:"name_en"`
	NameAbbEn string    `json:"name_abb_en"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoBankServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`banks.svc.info.start`)

	data, err := s.db.GetBankByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoBankServiceResponses{
		ID:        data.ID,
		NameTh:    data.NameTh,
		NameAbbTh: data.NameAbbTh,
		NameEn:    data.NameEn,
		NameAbbEn: data.NameAbbEn,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`banks.svc.info.success`)
	return resp, nil
}
