package payment_files

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoPaymentFileServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	PaymentID uuid.UUID `json:"payment_id"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoPaymentFileServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payment_files.svc.info.start`)

	var data *ent.PaymentFileEntity
	if isAdmin || memberID == uuid.Nil {
		file, err := s.db.GetPaymentFileByID(ctx, id)
		if err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return nil, err
		}
		data = file
	} else {
		file := new(ent.PaymentFileEntity)
		err := s.bunDB.DB().NewSelect().
			Model(file).
			Join("JOIN payments ON payments.id = payment_files.payment_id").
			Join("JOIN orders ON orders.payment_id = payments.id").
			Where("payment_files.id = ?", id).
			Where("orders.member_id = ?", memberID).
			Scan(ctx)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows
			}
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return nil, err
		}
		data = file
	}

	resp := &InfoPaymentFileServiceResponses{
		ID:        data.ID,
		PaymentID: data.PaymentID,
		FileID:    data.FileID,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`payment_files.svc.info.success`)
	return resp, nil
}
