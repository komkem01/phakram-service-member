package members

import (
	"context"
	"errors"
	"time"

	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (s *Service) ListMemberFilesService(ctx context.Context, req *entitiesdto.ListMemberFilesRequest) ([]*ent.MemberFileEntity, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.file.list.start`)

	data, page, err := s.file.ListMemberFiles(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	span.AddEvent(`members.svc.file.list.success`)
	return data, page, nil
}

func (s *Service) CreateMemberFileService(ctx context.Context, memberID uuid.UUID, fileID uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.file.create.start`)

	now := time.Now()
	memberFile := &ent.MemberFileEntity{
		ID:        uuid.New(),
		MemberID:  memberID,
		FileID:    fileID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(memberFile).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionCreated, ent.AuditActionCreated, "create_member_file", memberFile.ID, actionBy, "Created member file with ID "+memberFile.ID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionCreated, "create_member_file", memberFile.ID, actionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.file.create.success`)
	return nil
}

func (s *Service) InfoMemberFileService(ctx context.Context, memberID uuid.UUID, rowID uuid.UUID) (*ent.MemberFileEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.file.info.start`)

	data, err := s.file.GetMemberFileByID(ctx, rowID)
	if err != nil {
		return nil, err
	}
	if data.MemberID != memberID {
		return nil, errors.New("member file not found")
	}

	span.AddEvent(`members.svc.file.info.success`)
	return data, nil
}

func (s *Service) UpdateMemberFileService(ctx context.Context, memberID uuid.UUID, rowID uuid.UUID, fileID uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.file.update.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberFileEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", rowID).Where("deleted_at IS NULL").Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member file not found")
		}

		data.FileID = fileID
		data.UpdatedAt = now
		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionUpdated, ent.AuditActionUpdated, "update_member_file", data.ID, actionBy, "Updated member file with ID "+data.ID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionUpdated, "update_member_file", rowID, actionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.file.update.success`)
	return nil
}

func (s *Service) DeleteMemberFileService(ctx context.Context, memberID uuid.UUID, rowID uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.file.delete.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberFileEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", rowID).Where("deleted_at IS NULL").Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member file not found")
		}

		if _, err := tx.NewUpdate().Model(&ent.MemberFileEntity{}).Set("deleted_at = ?", now).Set("updated_at = ?", now).Where("id = ?", rowID).Exec(ctx); err != nil {
			return err
		}

		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionDeleted, ent.AuditActionDeleted, "delete_member_file", rowID, actionBy, "Deleted member file with ID "+rowID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionDeleted, "delete_member_file", rowID, actionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.file.delete.success`)
	return nil
}
