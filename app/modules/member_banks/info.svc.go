package member_banks

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoMemberBankServiceResponses struct {
	ID          uuid.UUID `json:"id"`
	MemberID    uuid.UUID `json:"member_id"`
	BankID      uuid.UUID `json:"bank_id"`
	BankNo      string    `json:"bank_no"`
	FirstnameTh string    `json:"firstname_th"`
	LastnameTh  string    `json:"lastname_th"`
	FirstnameEn string    `json:"firstname_en"`
	LastnameEn  string    `json:"lastname_en"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoMemberBankServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_banks.svc.info.start`)

	data, err := s.db.GetMemberBankByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoMemberBankServiceResponses{
		ID:          data.ID,
		MemberID:    data.MemberID,
		BankID:      data.BankID,
		BankNo:      data.BankNo,
		FirstnameTh: data.FirstnameTh,
		LastnameTh:  data.LastnameTh,
		FirstnameEn: data.FirstnameEn,
		LastnameEn:  data.LastnameEn,
		IsSystem:    data.IsSystem,
		IsActive:    data.IsActive,
		CreatedAt:   data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_banks.svc.info.success`)
	return resp, nil
}
