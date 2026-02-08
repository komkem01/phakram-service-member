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

var _ entitiesinf.CartItemEntity = (*Service)(nil)

func (s *Service) ListCartItems(ctx context.Context, req *entitiesdto.ListCartItemsRequest) ([]*ent.CartItemEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.CartItemEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"cart_id", "product_id"},
		[]string{"created_at", "cart_id"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Join("JOIN carts ON carts.id = cart_items.cart_id").
					Where("carts.member_id = ?", req.MemberID)
			}
			return q
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetCartItemByID(ctx context.Context, id uuid.UUID) (*ent.CartItemEntity, error) {
	data := new(ent.CartItemEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateCartItem(ctx context.Context, item *ent.CartItemEntity) error {
	_, err := s.db.NewInsert().
		Model(item).
		Exec(ctx)
	return err
}

func (s *Service) UpdateCartItem(ctx context.Context, item *ent.CartItemEntity) error {
	_, err := s.db.NewUpdate().
		Model(item).
		Where("id = ?", item.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteCartItem(ctx context.Context, itemID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.CartItemEntity{}).
		Where("id = ?", itemID).
		Exec(ctx)
	return err
}
