package banks

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateBankService struct {
	NameTh    string `json:"name_th"`
	NameAbbTh string `json:"name_abb_th"`
	NameEn    string `json:"name_en"`
	NameAbbEn string `json:"name_abb_en"`
}

func (s *Service) CreateBankService(ctx context.Context, req *CreateBankService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`banks.svc.create.start`)

	bank := &ent.BankEntity{
		ID:        uuid.New(),
		NameTh:    req.NameTh,
		NameAbbTh: req.NameAbbTh,
		NameEn:    req.NameEn,
		NameAbbEn: req.NameAbbEn,
	}
	if err := s.db.CreateBank(ctx, bank); err != nil {
		return err
	}
	span.AddEvent(`banks.svc.create.success`)
	return nil
}
