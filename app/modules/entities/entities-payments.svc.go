package entities

import (
	"context"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.PaymentEntity = (*Service)(nil)

func (s *Service) CreatePayment(ctx context.Context, payment *ent.PaymentEntity) error {
	_, err := s.db.NewInsert().
		Model(payment).
		Exec(ctx)
	return err
}

func (s *Service) GetPaymentByID(ctx context.Context, id uuid.UUID) (*ent.PaymentEntity, error) {
	data := new(ent.PaymentEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) UpdatePayment(ctx context.Context, payment *ent.PaymentEntity) error {
	_, err := s.db.NewUpdate().
		Model(payment).
		Where("id = ?", payment.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeletePayment(ctx context.Context, paymentID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.PaymentEntity{}).
		Where("id = ?", paymentID).
		Exec(ctx)
	return err
}
