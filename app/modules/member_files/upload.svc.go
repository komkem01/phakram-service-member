package member_files

import (
	"context"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UploadMemberFileService struct {
	MemberID uuid.UUID `json:"member_id"`
	FileName string    `json:"file_name"`
	FilePath string    `json:"file_path"`
	FileType string    `json:"file_type"`
	FileSize string    `json:"file_size"`
	MemberBy uuid.UUID `json:"member_by"`
}

type UploadMemberFileServiceResponse struct {
	FileID       uuid.UUID `json:"file_id"`
	MemberFileID uuid.UUID `json:"member_file_id"`
	FilePath     string    `json:"file_path"`
}

func (s *Service) UploadMemberFileService(ctx context.Context, req *UploadMemberFileService) (*UploadMemberFileServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_files.svc.upload.start`)

	var response *UploadMemberFileServiceResponse
	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		storageID := uuid.New()
		storage := &ent.StorageEntity{
			ID:            storageID,
			RefID:         req.MemberID,
			FileName:      req.FileName,
			FilePath:      req.FilePath,
			FileType:      req.FileType,
			FileSize:      req.FileSize,
			RelatedEntity: ent.RelateTypeMemberFile,
			UploadedBy:    req.MemberBy,
		}
		if _, err := tx.NewInsert().Model(storage).Exec(ctx); err != nil {
			return err
		}

		memberFile := &ent.MemberFileEntity{
			ID:       uuid.New(),
			MemberID: req.MemberID,
			FileID:   storageID,
		}
		if _, err := tx.NewInsert().Model(memberFile).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberBy
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "member_file",
			ActionID:     &memberFile.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Uploaded member file " + req.FileName,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		response = &UploadMemberFileServiceResponse{
			FileID:       storageID,
			MemberFileID: memberFile.ID,
			FilePath:     req.FilePath,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`member_files.svc.upload.success`)
	return response, nil
}
