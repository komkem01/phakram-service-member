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

var _ entitiesinf.MemberFileEntity = (*Service)(nil)

func (s *Service) ListMemberFiles(ctx context.Context, req *entitiesdto.ListMemberFilesRequest) ([]*ent.MemberFileEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.MemberFileEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"member_id", "file_id"},
		[]string{"created_at", "member_id", "file_id"},
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

func (s *Service) CreateMemberFile(ctx context.Context, memberFile *ent.MemberFileEntity) error {
	_, err := s.db.NewInsert().
		Model(memberFile).
		Exec(ctx)
	return err
}

func (s *Service) GetMemberFileByID(ctx context.Context, id uuid.UUID) (*ent.MemberFileEntity, error) {
	data := new(ent.MemberFileEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) UpdateMemberFile(ctx context.Context, memberFile *ent.MemberFileEntity) error {
	_, err := s.db.NewUpdate().
		Model(memberFile).
		Where("id = ?", memberFile.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberFile(ctx context.Context, memberFileID uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.MemberFileEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", memberFileID).
		Exec(ctx)
	return err
}
