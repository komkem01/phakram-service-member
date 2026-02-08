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

var _ entitiesinf.PaymentFileEntity = (*Service)(nil)

func (s *Service) ListPaymentFiles(ctx context.Context, req *entitiesdto.ListPaymentFilesRequest) ([]*ent.PaymentFileEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.PaymentFileEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"payment_id", "file_id"},
		[]string{"created_at", "payment_id"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Join("JOIN payments ON payments.id = payment_files.payment_id").
					Join("JOIN orders ON orders.payment_id = payments.id").
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

func (s *Service) CreatePaymentFile(ctx context.Context, file *ent.PaymentFileEntity) error {
	_, err := s.db.NewInsert().
		Model(file).
		Exec(ctx)
	return err
}

func (s *Service) UpdatePaymentFile(ctx context.Context, file *ent.PaymentFileEntity) error {
	_, err := s.db.NewUpdate().
		Model(file).
		Where("id = ?", file.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeletePaymentFile(ctx context.Context, fileID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.PaymentFileEntity{}).
		Where("id = ?", fileID).
		Exec(ctx)
	return err
}
