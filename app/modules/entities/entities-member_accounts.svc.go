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
		[]string{"email"},
		[]string{"created_at", "email"},
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

func (s *Service) GetMemberAccountByEmail(ctx context.Context, email string) (*ent.MemberAccountEntity, error) {
	data := new(ent.MemberAccountEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateMemberAccount(ctx context.Context, account *ent.MemberAccountEntity) error {
	_, err := s.db.NewInsert().
		Model(account).
		Exec(ctx)
	return err
}

func (s *Service) UpdateMemberAccount(ctx context.Context, account *ent.MemberAccountEntity) error {
	_, err := s.db.NewUpdate().
		Model(account).
		Where("id = ?", account.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberAccount(ctx context.Context, accountID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberAccountEntity{}).
		Where("id = ?", accountID).
		Exec(ctx)
	return err
}

func (s *Service) GetMemberAccountByMemberID(ctx context.Context, memberID uuid.UUID) (*ent.MemberAccountEntity, error) {
	data := new(ent.MemberAccountEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("member_id = ?", memberID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
