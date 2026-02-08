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

var _ entitiesinf.PaymentEntity = (*Service)(nil)

func (s *Service) ListPayments(ctx context.Context, req *entitiesdto.ListPaymentsRequest) ([]*ent.PaymentEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.PaymentEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"status"},
		[]string{"approved_at", "status"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Join("JOIN orders ON orders.payment_id = payments.id").
					Where("orders.member_id = ?", req.MemberID)
			}
			return q
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
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

func (s *Service) CreatePayment(ctx context.Context, payment *ent.PaymentEntity) error {
	_, err := s.db.NewInsert().
		Model(payment).
		Exec(ctx)
	return err
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
