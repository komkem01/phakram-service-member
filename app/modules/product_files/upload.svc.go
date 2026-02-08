package product_files

import (
	"context"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UploadProductFileService struct {
	ProductID uuid.UUID `json:"product_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	FileSize  string    `json:"file_size"`
	MemberID  uuid.UUID `json:"member_id"`
}

type UploadProductFileServiceResponse struct {
	FileID        uuid.UUID `json:"file_id"`
	ProductFileID uuid.UUID `json:"product_file_id"`
	FilePath      string    `json:"file_path"`
}

func (s *Service) UploadProductFileService(ctx context.Context, req *UploadProductFileService) (*UploadProductFileServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_files.svc.upload.start`)

	var response *UploadProductFileServiceResponse
	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		storageID := uuid.New()
		storage := &ent.StorageEntity{
			ID:            storageID,
			RefID:         req.ProductID,
			FileName:      req.FileName,
			FilePath:      req.FilePath,
			FileType:      req.FileType,
			FileSize:      req.FileSize,
			RelatedEntity: ent.RelateTypeProductFile,
			UploadedBy:    req.MemberID,
		}
		if _, err := tx.NewInsert().Model(storage).Exec(ctx); err != nil {
			return err
		}

		productFile := &ent.ProductFileEntity{
			ID:        uuid.New(),
			ProductID: req.ProductID,
			FileID:    storageID,
		}
		if _, err := tx.NewInsert().Model(productFile).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "product_file",
			ActionID:     &productFile.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Uploaded product file " + req.FileName,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		response = &UploadProductFileServiceResponse{
			FileID:        storageID,
			ProductFileID: productFile.ID,
			FilePath:      req.FilePath,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`product_files.svc.upload.success`)
	return response, nil
}
