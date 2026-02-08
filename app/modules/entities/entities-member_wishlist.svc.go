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

var _ entitiesinf.MemberWishlistEntity = (*Service)(nil)

func (s *Service) ListMemberWishlist(ctx context.Context, req *entitiesdto.ListMemberWishlistRequest) ([]*ent.MemberWishlistEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberWishlistEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_id", "product_id"},
		[]string{"created_at", "member_id", "product_id"},
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

func (s *Service) GetMemberWishlistByID(ctx context.Context, id uuid.UUID) (*ent.MemberWishlistEntity, error) {
	data := new(ent.MemberWishlistEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateMemberWishlist(ctx context.Context, wishlist *ent.MemberWishlistEntity) error {
	_, err := s.db.NewInsert().
		Model(wishlist).
		Exec(ctx)
	return err
}

func (s *Service) UpdateMemberWishlist(ctx context.Context, wishlist *ent.MemberWishlistEntity) error {
	_, err := s.db.NewUpdate().
		Model(wishlist).
		Where("id = ?", wishlist.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberWishlist(ctx context.Context, wishlistID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberWishlistEntity{}).
		Where("id = ?", wishlistID).
		Exec(ctx)
	return err
}
