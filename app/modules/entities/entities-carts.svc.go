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

var _ entitiesinf.CartEntity = (*Service)(nil)

func (s *Service) ListCarts(ctx context.Context, req *entitiesdto.ListCartsRequest) ([]*ent.CartEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.CartEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_id"},
		[]string{"created_at", "member_id"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Where("member_id = ?", req.MemberID)
			}
			return q
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetCartByID(ctx context.Context, id uuid.UUID) (*ent.CartEntity, error) {
	data := new(ent.CartEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateCart(ctx context.Context, cart *ent.CartEntity) error {
	_, err := s.db.NewInsert().
		Model(cart).
		Exec(ctx)
	return err
}

func (s *Service) UpdateCart(ctx context.Context, cart *ent.CartEntity) error {
	_, err := s.db.NewUpdate().
		Model(cart).
		Where("id = ?", cart.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.CartEntity{}).
		Where("id = ?", cartID).
		Exec(ctx)
	return err
}
