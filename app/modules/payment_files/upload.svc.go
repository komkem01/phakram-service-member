package payment_files

import (
	"context"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UploadPaymentFileService struct {
	PaymentID uuid.UUID `json:"payment_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	FileSize  string    `json:"file_size"`
	MemberID  uuid.UUID `json:"member_id"`
}

type UploadPaymentFileServiceResponse struct {
	FileID        uuid.UUID `json:"file_id"`
	PaymentFileID uuid.UUID `json:"payment_file_id"`
	FilePath      string    `json:"file_path"`
}

func (s *Service) UploadPaymentFileService(ctx context.Context, req *UploadPaymentFileService) (*UploadPaymentFileServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payment_files.svc.upload.start`)

	var response *UploadPaymentFileServiceResponse
	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		storageID := uuid.New()
		storage := &ent.StorageEntity{
			ID:            storageID,
			RefID:         req.PaymentID,
			FileName:      req.FileName,
			FilePath:      req.FilePath,
			FileType:      req.FileType,
			FileSize:      req.FileSize,
			RelatedEntity: ent.RelateTypePaymentFile,
			UploadedBy:    req.MemberID,
		}
		if _, err := tx.NewInsert().Model(storage).Exec(ctx); err != nil {
			return err
		}

		paymentFile := &ent.PaymentFileEntity{
			ID:        uuid.New(),
			PaymentID: req.PaymentID,
			FileID:    storageID,
		}
		if _, err := tx.NewInsert().Model(paymentFile).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "payment_file",
			ActionID:     &paymentFile.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Uploaded payment file " + req.FileName,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		response = &UploadPaymentFileServiceResponse{
			FileID:        storageID,
			PaymentFileID: paymentFile.ID,
			FilePath:      req.FilePath,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`payment_files.svc.upload.success`)
	return response, nil
}
