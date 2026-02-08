package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.OrderEntity = (*Service)(nil)

func (s *Service) ListOrders(ctx context.Context, req *entitiesdto.ListOrdersRequest) ([]*ent.OrderEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.OrderEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"order_no", "status"},
		[]string{"created_at", "order_no"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Where("member_id = ?", req.MemberID)
			}
			if req.Status != "" {
				q.Where("status = ?", req.Status)
			}
			if req.StartDate > 0 {
				q.Where("created_at >= ?", unixToTime(req.StartDate))
			}
			if req.EndDate > 0 {
				q.Where("created_at <= ?", unixToTime(req.EndDate))
			}
			if req.Search != "" {
				like := "%" + req.Search + "%"
				q.Join("JOIN members AS m ON m.id = orders.member_id").
					Where("orders.order_no ILIKE ? OR m.firstname_th ILIKE ? OR m.lastname_th ILIKE ? OR m.firstname_en ILIKE ? OR m.lastname_en ILIKE ?", like, like, like, like, like)
			}
			return q
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func unixToTime(value int64) time.Time {
	if value <= 0 {
		return time.Time{}
	}
	if value > 1_000_000_000_000 {
		return time.Unix(0, value*int64(time.Millisecond))
	}
	return time.Unix(value, 0)
}

func (s *Service) GetOrderByID(ctx context.Context, id uuid.UUID) (*ent.OrderEntity, error) {
	data := new(ent.OrderEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateOrder(ctx context.Context, order *ent.OrderEntity) error {
	_, err := s.db.NewInsert().
		Model(order).
		Exec(ctx)
	return err
}

func (s *Service) UpdateOrder(ctx context.Context, order *ent.OrderEntity) error {
	_, err := s.db.NewUpdate().
		Model(order).
		Where("id = ?", order.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.OrderEntity{}).
		Where("id = ?", orderID).
		Exec(ctx)
	return err
}
