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

var _ entitiesinf.MemberPaymentEntity = (*Service)(nil)

func (s *Service) ListMemberPayments(ctx context.Context, req *entitiesdto.ListMemberPaymentsRequest) ([]*ent.MemberPaymentEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberPaymentEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_id", "payment_id"},
		[]string{"created_at", "member_id", "payment_id"},
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

func (s *Service) CreateMemberPayment(ctx context.Context, memberPayment *ent.MemberPaymentEntity) error {
	_, err := s.db.NewInsert().
		Model(memberPayment).
		Exec(ctx)
	return err
}

func (s *Service) GetMemberPaymentByID(ctx context.Context, id uuid.UUID) (*ent.MemberPaymentEntity, error) {
	data := new(ent.MemberPaymentEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) UpdateMemberPayment(ctx context.Context, memberPayment *ent.MemberPaymentEntity) error {
	_, err := s.db.NewUpdate().
		Model(memberPayment).
		Where("id = ?", memberPayment.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberPayment(ctx context.Context, memberPaymentID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberPaymentEntity{}).
		Where("id = ?", memberPaymentID).
		Exec(ctx)
	return err
}
