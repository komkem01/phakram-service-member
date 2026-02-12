package entities

import (
	"context"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberBankEntity = (*Service)(nil)

func (s *Service) CreateMemberBank(ctx context.Context, memberBank *ent.MemberBankEntity) error {
	data := ent.MemberBankEntity{
		ID:          memberBank.ID,
		MemberID:    memberBank.MemberID,
		BankID:      memberBank.BankID,
		BankNo:      memberBank.BankNo,
		FirstnameTh: memberBank.FirstnameTh,
		LastnameTh:  memberBank.LastnameTh,
		FirstnameEn: memberBank.FirstnameEn,
		LastnameEn:  memberBank.LastnameEn,
		IsDefault:   memberBank.IsDefault,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := s.db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func (s *Service) GetMemberBankByID(ctx context.Context, id uuid.UUID) (*ent.MemberBankEntity, error) {
	var memberBank ent.MemberBankEntity
	err := s.db.NewSelect().Model(&memberBank).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &memberBank, nil
}

func (s *Service) UpdateMemberBank(ctx context.Context, memberBank *ent.MemberBankEntity) error {
	_, err := s.db.NewUpdate().Model(memberBank).Where("id = ?", memberBank.ID).Exec(ctx)
	return err
}

func (s *Service) DeleteMemberBank(ctx context.Context, memberBankID uuid.UUID) error {
	_, err := s.db.NewDelete().Model(&ent.MemberBankEntity{}).Where("id = ?", memberBankID).Exec(ctx)
	return err
}
