package entities

import (
	"context"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberPaymentEntity = (*Service)(nil)

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
