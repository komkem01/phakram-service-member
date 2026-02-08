package entities

import (
	"context"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberAddressEntity = (*Service)(nil)

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
