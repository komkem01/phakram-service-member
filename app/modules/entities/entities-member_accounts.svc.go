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

var _ entitiesinf.MemberAccountEntity = (*Service)(nil)

func (s *Service) ListMemberAccounts(ctx context.Context, req *entitiesdto.ListMemberAccountsRequest) ([]*ent.MemberAccountEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberAccountEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"email", "member_id"},
		[]string{"created_at", "email", "member_id"},
		func(selQ *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				selQ.Where("member_id = ?", req.MemberID)
			}
			return selQ
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) CreateMemberAccount(ctx context.Context, memberAccount *ent.MemberAccountEntity) error {
	_, err := s.db.NewInsert().
		Model(memberAccount).
		Exec(ctx)
	return err
}

func (s *Service) GetMemberAccountByID(ctx context.Context, id uuid.UUID) (*ent.MemberAccountEntity, error) {
	data := new(ent.MemberAccountEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) UpdateMemberAccount(ctx context.Context, memberAccount *ent.MemberAccountEntity) error {
	_, err := s.db.NewUpdate().
		Model(memberAccount).
		Where("id = ?", memberAccount.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberAccount(ctx context.Context, memberAccountID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberAccountEntity{}).
		Where("id = ?", memberAccountID).
		Exec(ctx)
	return err
}
