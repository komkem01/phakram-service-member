package member_banks

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdateMemberBankService struct {
	MemberID    uuid.UUID `json:"member_id"`
	BankID      uuid.UUID `json:"bank_id"`
	BankNo      string    `json:"bank_no"`
	FirstnameTh string    `json:"firstname_th"`
	LastnameTh  string    `json:"lastname_th"`
	FirstnameEn string    `json:"firstname_en"`
	LastnameEn  string    `json:"lastname_en"`
	IsSystem    *bool     `json:"is_system"`
	IsActive    *bool     `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateMemberBankService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_banks.svc.update.start`)

	data, err := s.db.GetMemberBankByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.MemberID != uuid.Nil {
		data.MemberID = req.MemberID
	}
	if req.BankID != uuid.Nil {
		data.BankID = req.BankID
	}
	if req.BankNo != "" {
		data.BankNo = req.BankNo
	}
	if req.FirstnameTh != "" {
		data.FirstnameTh = req.FirstnameTh
	}
	if req.LastnameTh != "" {
		data.LastnameTh = req.LastnameTh
	}
	if req.FirstnameEn != "" {
		data.FirstnameEn = req.FirstnameEn
	}
	if req.LastnameEn != "" {
		data.LastnameEn = req.LastnameEn
	}
	if req.IsSystem != nil {
		data.IsSystem = *req.IsSystem
	}
	if req.IsActive != nil {
		data.IsActive = *req.IsActive
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditUpdate,
			ActionType:   "member_bank",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated member bank " + data.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`member_banks.svc.update.success`)
	return nil
}
