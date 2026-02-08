package member_banks

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateMemberBankService struct {
	MemberID    uuid.UUID `json:"member_id"`
	BankID      uuid.UUID `json:"bank_id"`
	BankNo      string    `json:"bank_no"`
	FirstnameTh string    `json:"firstname_th"`
	LastnameTh  string    `json:"lastname_th"`
	FirstnameEn string    `json:"firstname_en"`
	LastnameEn  string    `json:"lastname_en"`
	IsSystem    bool      `json:"is_system"`
	IsActive    *bool     `json:"is_active"`
}

func (s *Service) CreateMemberBankService(ctx context.Context, req *CreateMemberBankService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_banks.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
		bank := &ent.MemberBankEntity{
			ID:          uuid.New(),
			MemberID:    req.MemberID,
			BankID:      req.BankID,
			BankNo:      req.BankNo,
			FirstnameTh: req.FirstnameTh,
			LastnameTh:  req.LastnameTh,
			FirstnameEn: req.FirstnameEn,
			LastnameEn:  req.LastnameEn,
			IsSystem:    req.IsSystem,
			IsActive:    isActive,
		}
		if _, err := tx.NewInsert().Model(bank).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "member_bank",
			ActionID:     &bank.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created member bank " + bank.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	span.AddEvent(`member_banks.svc.create.success`)
	return nil
}
