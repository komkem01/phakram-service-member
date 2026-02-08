package banks

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateBankService struct {
	NameTh    string `json:"name_th"`
	NameAbbTh string `json:"name_abb_th"`
	NameEn    string `json:"name_en"`
	NameAbbEn string `json:"name_abb_en"`
	IsActive  bool   `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateBankService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`banks.svc.update.start`)

	data, err := s.db.GetBankByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.NameTh != "" {
		data.NameTh = req.NameTh
	}
	if req.NameAbbTh != "" {
		data.NameAbbTh = req.NameAbbTh
	}
	if req.NameEn != "" {
		data.NameEn = req.NameEn
	}
	if req.NameAbbEn != "" {
		data.NameAbbEn = req.NameAbbEn
	}

	if err := s.db.UpdateBank(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`banks.svc.update.success`)
	return nil
}
