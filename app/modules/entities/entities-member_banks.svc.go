package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberBankEntity = (*Service)(nil)

func (s *Service) ListMemberBanks(ctx context.Context, req *entitiesdto.ListMemberBanksRequest) ([]*ent.MemberBankEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberBankEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"bank_no", "firstname_th", "lastname_th", "firstname_en", "lastname_en"},
		[]string{"created_at", "bank_no", "firstname_th", "lastname_th"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Where("member_id = ?", req.MemberID)
			}
			return q
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetMemberBankByID(ctx context.Context, id uuid.UUID) (*ent.MemberBankEntity, error) {
	data := new(ent.MemberBankEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateMemberBank(ctx context.Context, bank *ent.MemberBankEntity) error {
	_, err := s.db.NewInsert().
		Model(bank).
		Exec(ctx)
	return err
}

func (s *Service) UpdateMemberBank(ctx context.Context, bank *ent.MemberBankEntity) error {
	_, err := s.db.NewUpdate().
		Model(bank).
		Where("id = ?", bank.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberBank(ctx context.Context, bankID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberBankEntity{}).
		Where("id = ?", bankID).
		Exec(ctx)
	return err
}
