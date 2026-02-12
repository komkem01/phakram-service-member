package entities

import (
	"context"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.PaymentFileEntity = (*Service)(nil)

func (s *Service) CreatePaymentFile(ctx context.Context, paymentFile *ent.PaymentFileEntity) error {
	_, err := s.db.NewInsert().
		Model(paymentFile).
		Exec(ctx)
	return err
}

func (s *Service) GetPaymentFileByID(ctx context.Context, id uuid.UUID) (*ent.PaymentFileEntity, error) {
	data := new(ent.PaymentFileEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) UpdatePaymentFile(ctx context.Context, paymentFile *ent.PaymentFileEntity) error {
	_, err := s.db.NewUpdate().
		Model(paymentFile).
		Where("id = ?", paymentFile.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeletePaymentFile(ctx context.Context, paymentFileID uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.PaymentFileEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", paymentFileID).
		Exec(ctx)
	return err
}
