package payment_files

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListPaymentFileServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListPaymentFileServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	PaymentID uuid.UUID `json:"payment_id"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListPaymentFileServiceRequest) ([]*ListPaymentFileServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payment_files.svc.list.start`)

	data, page, err := s.db.ListPaymentFiles(ctx, &entitiesdto.ListPaymentFilesRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListPaymentFileServiceResponses
	for _, item := range data {
		temp := &ListPaymentFileServiceResponses{
			ID:        item.ID,
			PaymentID: item.PaymentID,
			FileID:    item.FileID,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`payment_files.svc.list.copy`)
	return response, page, nil
}
