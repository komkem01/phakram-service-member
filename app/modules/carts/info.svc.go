package carts

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoCartServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoCartServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.info.start`)

	data, err := s.db.GetCartByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoCartServiceResponses{
		ID:        data.ID,
		MemberID:  data.MemberID,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`carts.svc.info.success`)
	return resp, nil
}
