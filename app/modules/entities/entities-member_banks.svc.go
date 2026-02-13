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

var _ entitiesinf.MemberBankEntity = (*Service)(nil)

func (s *Service) ListMemberBanks(ctx context.Context, req *entitiesdto.ListMemberBanksRequest) ([]*ent.MemberBankEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberBankEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_id", "bank_no", "firstname_th", "lastname_th"},
		[]string{"created_at", "member_id", "bank_no", "firstname_th", "lastname_th"},
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
