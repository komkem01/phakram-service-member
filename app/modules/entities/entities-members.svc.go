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

var _ entitiesinf.MemberEntity = (*Service)(nil)

func (s *Service) CreateMember(ctx context.Context, member *ent.MemberEntity) error {
	data := ent.MemberEntity{
		MemberNo:     member.MemberNo,
		PrefixID:     member.PrefixID,
		GenderID:     member.GenderID,
		FirstnameTh:  member.FirstnameTh,
		LastnameTh:   member.LastnameTh,
		FirstnameEn:  member.FirstnameEn,
		LastnameEn:   member.LastnameEn,
		Role:         member.Role,
		Phone:        member.Phone,
		Registration: member.Registration,
		CreatedAt:    time.Now(),
	}
	_, err := s.db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func (s *Service) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error) {
	var member ent.MemberEntity
	err := s.db.NewSelect().Model(&member).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *Service) ListMembers(ctx context.Context, req *entitiesdto.ListMembersRequest) ([]*ent.MemberEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_no", "firstname_th", "lastname_th", "firstname_en", "lastname_en", "phone"},
		[]string{"created_at", "member_no", "firstname_th", "lastname_th", "firstname_en", "lastname_en"},
		func(selQ *bun.SelectQuery) *bun.SelectQuery {
			if req.Role != "" {
				selQ = selQ.Where("role = ?", req.Role)
			}
			return selQ
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) UpdateMember(ctx context.Context, member *ent.MemberEntity) error {
	_, err := s.db.NewUpdate().Model(member).Where("id = ?", member.ID).Exec(ctx)
	return err
}

func (s *Service) DeleteMember(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.MemberEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

// admin service
func (s *Service) CreateAdminMember(ctx context.Context, member *ent.MemberEntity) error {
	data := ent.MemberEntity{
		MemberNo:      member.MemberNo,
		TierID:        member.TierID,
		StatusID:      member.StatusID,
		PrefixID:      member.PrefixID,
		GenderID:      member.GenderID,
		FirstnameTh:   member.FirstnameTh,
		LastnameTh:    member.LastnameTh,
		FirstnameEn:   member.FirstnameEn,
		LastnameEn:    member.LastnameEn,
		Role:          member.Role,
		Phone:         member.Phone,
		TotalSpent:    member.TotalSpent,
		CurrentPoints: member.CurrentPoints,
		Registration:  member.Registration,
		CreatedAt:     time.Now(),
	}
	_, err := s.db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func (s *Service) GetAdminMemberByID(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error) {
	var member ent.MemberEntity
	err := s.db.NewSelect().Model(&member).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *Service) UpdateAdminMember(ctx context.Context, member *ent.MemberEntity) error {
	_, err := s.db.NewUpdate().Model(member).Where("id = ?", member.ID).Exec(ctx)
	return err
}

func (s *Service) DeleteAdminMember(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.MemberEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

// member service by admin
func (s *Service) CreateMemberByAdmin(ctx context.Context, member *ent.MemberEntity) error {
	data := ent.MemberEntity{
		MemberNo:      member.MemberNo,
		TierID:        member.TierID,
		StatusID:      member.StatusID,
		PrefixID:      member.PrefixID,
		GenderID:      member.GenderID,
		FirstnameTh:   member.FirstnameTh,
		LastnameTh:    member.LastnameTh,
		FirstnameEn:   member.FirstnameEn,
		LastnameEn:    member.LastnameEn,
		Role:          member.Role,
		Phone:         member.Phone,
		TotalSpent:    member.TotalSpent,
		CurrentPoints: member.CurrentPoints,
		Registration:  member.Registration,
		CreatedAt:     time.Now(),
	}
	_, err := s.db.NewInsert().Model(&data).Exec(ctx)
	return err
}

func (s *Service) GetMemberByIDByAdmin(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error) {
	var member ent.MemberEntity
	err := s.db.NewSelect().Model(&member).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *Service) UpdateMemberByAdmin(ctx context.Context, member *ent.MemberEntity) error {
	_, err := s.db.NewUpdate().Model(member).Where("id = ?", member.ID).Exec(ctx)
	return err
}

func (s *Service) DeleteMemberByAdmin(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.MemberEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
