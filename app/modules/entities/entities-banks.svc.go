package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.BankEntity = (*Service)(nil)

func (s *Service) ListBanks(ctx context.Context, req *entitiesdto.ListBanksRequest) ([]*ent.BankEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.BankEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name_th", "name_abb_th", "name_en", "name_abb_en"},
		[]string{"created_at", "name_th", "name_abb_th", "name_en", "name_abb_en"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetBankByID(ctx context.Context, id uuid.UUID) (*ent.BankEntity, error) {
	data := new(ent.BankEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateBank(ctx context.Context, bank *ent.BankEntity) error {
	_, err := s.db.NewInsert().
		Model(bank).
		Exec(ctx)
	return err
}

func (s *Service) UpdateBank(ctx context.Context, bank *ent.BankEntity) error {
	_, err := s.db.NewUpdate().
		Model(bank).
		Where("id = ?", bank.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteBank(ctx context.Context, bankID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.BankEntity{}).
		Where("id = ?", bankID).
		Exec(ctx)
	return err
}
