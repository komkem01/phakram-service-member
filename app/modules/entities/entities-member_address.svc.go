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

var _ entitiesinf.MemberAddressEntity = (*Service)(nil)

func (s *Service) ListMemberAddresses(ctx context.Context, req *entitiesdto.ListMemberAddressesRequest) ([]*ent.MemberAddressEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberAddressEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_id", "phone", "first_name", "last_name"},
		[]string{"created_at", "member_id", "phone", "first_name", "last_name"},
		func(selQ *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				selQ.Where("member_id = ?", req.MemberID)
			}
			return selQ
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return data, page, nil
}

func (s *Service) CreateMemberAddress(ctx context.Context, memberAddress *ent.MemberAddressEntity) error {
	data := ent.MemberAddressEntity{
		ID:            memberAddress.ID,
		MemberID:      memberAddress.MemberID,
		FirstName:     memberAddress.FirstName,
		LastName:      memberAddress.LastName,
		Phone:         memberAddress.Phone,
		AddressNo:     memberAddress.AddressNo,
		Village:       memberAddress.Village,
		Alley:         memberAddress.Alley,
		SubDistrictID: memberAddress.SubDistrictID,
		DistrictID:    memberAddress.DistrictID,
		ProvinceID:    memberAddress.ProvinceID,
		ZipcodeID:     memberAddress.ZipcodeID,
		IsDefault:     memberAddress.IsDefault,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, err := s.db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func (s *Service) GetMemberAddressByID(ctx context.Context, id uuid.UUID) (*ent.MemberAddressEntity, error) {
	var memberAddress ent.MemberAddressEntity
	err := s.db.NewSelect().Model(&memberAddress).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &memberAddress, nil
}

func (s *Service) UpdateMemberAddress(ctx context.Context, memberAddress *ent.MemberAddressEntity) error {
	_, err := s.db.NewUpdate().Model(memberAddress).Where("id = ?", memberAddress.ID).Exec(ctx)
	return err
}

func (s *Service) DeleteMemberAddress(ctx context.Context, memberAddressID uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.MemberAddressEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", memberAddressID).
		Exec(ctx)
	return err
}
