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

var _ entitiesinf.MemberEntity = (*Service)(nil)

func (s *Service) ListMembers(ctx context.Context, req *entitiesdto.ListMembersRequest) ([]*ent.MemberEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_no", "firstname_th", "lastname_th", "firstname_en", "lastname_en", "phone"},
		[]string{"created_at", "member_no", "firstname_th", "lastname_th", "phone"},
		func(q *bun.SelectQuery) *bun.SelectQuery {
			if req.MemberID != uuid.Nil {
				q.Where("id = ?", req.MemberID)
			}
			return q
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error) {
	data := new(ent.MemberEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetMemberByPhone(ctx context.Context, phone string) (*ent.MemberEntity, error) {
	data := new(ent.MemberEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("phone = ?", phone).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateMember(ctx context.Context, member *ent.MemberEntity) error {
	_, err := s.db.NewInsert().
		Model(member).
		Exec(ctx)
	return err
}

func (s *Service) UpdateMember(ctx context.Context, member *ent.MemberEntity) error {
	_, err := s.db.NewUpdate().
		Model(member).
		Where("id = ?", member.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.MemberEntity{}).
		Where("id = ?", memberID).
		Exec(ctx)
	return err
}
