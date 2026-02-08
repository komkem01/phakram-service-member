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

var _ entitiesinf.MemberAddressEntity = (*Service)(nil)

func (s *Service) ListMemberAddresses(ctx context.Context, req *entitiesdto.ListMemberAddressesRequest) ([]*ent.MemberAddressEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberAddressEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"first_name", "last_name", "phone", "address_no", "village"},
		[]string{"created_at", "first_name", "last_name", "phone"},
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

func (s *Service) GetMemberAddressByID(ctx context.Context, id uuid.UUID) (*ent.MemberAddressEntity, error) {
	data := new(ent.MemberAddressEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateMemberAddress(ctx context.Context, address *ent.MemberAddressEntity) error {
	_, err := s.db.NewInsert().
		Model(address).
		Exec(ctx)
	return err
}

func (s *Service) UpdateMemberAddress(ctx context.Context, address *ent.MemberAddressEntity) error {
	_, err := s.db.NewUpdate().
		Model(address).
		Where("id = ?", address.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberAddress(ctx context.Context, addressID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberAddressEntity{}).
		Where("id = ?", addressID).
		Exec(ctx)
	return err
}
