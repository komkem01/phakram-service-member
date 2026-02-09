package tiers

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateTierService struct {
	NameTh       string          `json:"name_th"`
	NameEn       string          `json:"name_en"`
	MinSpending  decimal.Decimal `json:"min_spending"`
	IsActive     bool            `json:"is_active"`
	DiscountRate decimal.Decimal `json:"discount_rate"`
}

func (s *Service) CreateTierService(ctx context.Context, req *CreateTierService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`tiers.svc.create.start`)

	id := uuid.New()
	tier := &ent.TierEntity{
		ID:           id,
		NameTh:       req.NameTh,
		NameEn:       req.NameEn,
		MinSpending:  req.MinSpending,
		IsActive:     req.IsActive,
		DiscountRate: req.DiscountRate,
	}
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(tier).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_tier",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created tier with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`tiers.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_tier",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create tier failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`tiers.svc.create.success`)
	return nil
}
