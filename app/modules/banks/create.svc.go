package banks

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateBankService struct {
	NameTh    string `json:"name_th"`
	NameAbbTh string `json:"name_abb_th"`
	NameEn    string `json:"name_en"`
	NameAbbEn string `json:"name_abb_en"`
	IsActive  bool   `json:"is_active"`
}

func (s *Service) CreateBankService(ctx context.Context, req *CreateBankService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`banks.svc.create.start`)

	id := uuid.New()
	bank := &ent.BankEntity{
		ID:        id,
		NameTh:    req.NameTh,
		NameAbbTh: req.NameAbbTh,
		NameEn:    req.NameEn,
		NameAbbEn: req.NameAbbEn,
		IsActive:  req.IsActive,
	}
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(bank).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_bank",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created bank with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`banks.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_bank",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create bank failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`banks.svc.create.success`)
	return nil
}
