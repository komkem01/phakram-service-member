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

var _ entitiesinf.MemberTransactionEntity = (*Service)(nil)

func (s *Service) ListMemberTransactions(ctx context.Context, req *entitiesdto.ListMemberTransactionsRequest) ([]*ent.MemberTransactionEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberTransactionEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"action", "details"},
		[]string{"created_at", "action"},
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

func (s *Service) GetMemberTransactionByID(ctx context.Context, id uuid.UUID) (*ent.MemberTransactionEntity, error) {
	data := new(ent.MemberTransactionEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateMemberTransaction(ctx context.Context, transaction *ent.MemberTransactionEntity) error {
	_, err := s.db.NewInsert().
		Model(transaction).
		Exec(ctx)
	return err
}

func (s *Service) UpdateMemberTransaction(ctx context.Context, transaction *ent.MemberTransactionEntity) error {
	_, err := s.db.NewUpdate().
		Model(transaction).
		Where("id = ?", transaction.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberTransaction(ctx context.Context, transactionID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberTransactionEntity{}).
		Where("id = ?", transactionID).
		Exec(ctx)
	return err
}
