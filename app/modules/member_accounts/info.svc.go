package member_accounts

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoMemberAccountServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoMemberAccountServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_accounts.svc.info.start`)

	data, err := s.db.GetMemberAccountByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoMemberAccountServiceResponses{
		ID:        data.ID,
		MemberID:  data.MemberID,
		Email:     data.Email,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_accounts.svc.info.success`)
	return resp, nil
}

func (s *Service) InfoByEmailService(ctx context.Context, email string, memberID uuid.UUID, isAdmin bool) (*InfoMemberAccountServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_accounts.svc.info_by_email.start`)

	data, err := s.db.GetMemberAccountByEmail(ctx, email)
	if err != nil {
		log.With(slog.Any(`email`, email)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoMemberAccountServiceResponses{
		ID:        data.ID,
		MemberID:  data.MemberID,
		Email:     data.Email,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_accounts.svc.info_by_email.success`)
	return resp, nil
}

func (s *Service) InfoByMemberIDService(ctx context.Context, memberID uuid.UUID, requesterMemberID uuid.UUID, isAdmin bool) (*InfoMemberAccountServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_accounts.svc.info_by_member_id.start`)
	data, err := s.db.GetMemberAccountByMemberID(ctx, memberID)
	if err != nil {
		log.With(slog.Any(`member_id`, memberID)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && requesterMemberID != uuid.Nil && data.MemberID != requesterMemberID {
		return nil, sql.ErrNoRows
	}
	resp := &InfoMemberAccountServiceResponses{
		ID:        data.ID,
		MemberID:  data.MemberID,
		Email:     data.Email,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_accounts.svc.info_by_member_id.success`)
	return resp, nil
}
