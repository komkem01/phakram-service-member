package member_transactions

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type UpdateMemberTransactionService struct {
	MemberID uuid.UUID `json:"member_id"`
	Action   string    `json:"action"`
	Details  string    `json:"details"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateMemberTransactionService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_transactions.svc.update.start`)

	data, err := s.db.GetMemberTransactionByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.MemberID != uuid.Nil {
		data.MemberID = req.MemberID
	}
	if req.Action != "" {
		data.Action = ent.ActionTypeEnum(req.Action)
	}
	if req.Details != "" {
		data.Details = req.Details
	}

	if err := s.db.UpdateMemberTransaction(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`member_transactions.svc.update.success`)
	return nil
}
