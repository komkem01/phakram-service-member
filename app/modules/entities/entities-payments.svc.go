package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.PaymentEntity = (*Service)(nil)

func (s *Service) ListPayments(ctx context.Context, req *entitiesdto.ListPaymentsRequest) ([]*ent.PaymentEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.PaymentEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"status"},
		[]string{"id", "amount", "status", "approved_at"},
		func(selQ *bun.SelectQuery) *bun.SelectQuery {
			if strings.TrimSpace(req.Status) != "" {
				selQ.Where("status = ?", req.Status)
			}
			return selQ
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return data, page, nil
}

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
