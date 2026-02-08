package member_transactions

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoMemberTransactionServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	CreatedAt string    `json:"created_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoMemberTransactionServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_transactions.svc.info.start`)

	data, err := s.db.GetMemberTransactionByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoMemberTransactionServiceResponses{
		ID:        data.ID,
		MemberID:  data.MemberID,
		Action:    string(data.Action),
		Details:   data.Details,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_transactions.svc.info.success`)
	return resp, nil
}
