package entities

import (
	"context"
	"os"
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
	loc := orderTimeLocation()

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"order_no", "member_id", "status"},
		[]string{"created_at", "order_no", "member_id", "status"},
		func(selQ *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				selQ.Where("member_id = ?", req.MemberID)
			}
			if req.Status != "" {
				selQ.Where("status = ?", req.Status)
			}
			if req.Search != "" {
				selQ.Where("order_no ILIKE ?", "%"+req.Search+"%")
			}
			if req.StartDate > 0 {
				selQ.Where("created_at >= ?", time.Unix(req.StartDate, 0).In(loc))
			}
			if req.EndDate > 0 {
				selQ.Where("created_at <= ?", time.Unix(req.EndDate, 0).In(loc))
			}
			return selQ
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
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

func orderTimeLocation() *time.Location {
	keys := []string{"DATABASE_SQL__TIME_ZONE", "DATABASE_SQL__TIMEZONE", "DB_TIMEZONE"}
	for _, key := range keys {
		if tz := os.Getenv(key); tz != "" {
			loc, err := time.LoadLocation(tz)
			if err == nil {
				return loc
			}
			break
		}
	}
	return time.Local
}
