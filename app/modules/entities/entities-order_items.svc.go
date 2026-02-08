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

var _ entitiesinf.OrderItemEntity = (*Service)(nil)

func (s *Service) ListOrderItems(ctx context.Context, req *entitiesdto.ListOrderItemsRequest) ([]*ent.OrderItemEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.OrderItemEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"order_id", "product_id"},
		[]string{"created_at", "order_id"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Join("JOIN orders ON orders.id = order_items.order_id").
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

func (s *Service) GetOrderItemByID(ctx context.Context, id uuid.UUID) (*ent.OrderItemEntity, error) {
	data := new(ent.OrderItemEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateOrderItem(ctx context.Context, item *ent.OrderItemEntity) error {
	_, err := s.db.NewInsert().
		Model(item).
		Exec(ctx)
	return err
}

func (s *Service) UpdateOrderItem(ctx context.Context, item *ent.OrderItemEntity) error {
	_, err := s.db.NewUpdate().
		Model(item).
		Where("id = ?", item.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteOrderItem(ctx context.Context, itemID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.OrderItemEntity{}).
		Where("id = ?", itemID).
		Exec(ctx)
	return err
}
