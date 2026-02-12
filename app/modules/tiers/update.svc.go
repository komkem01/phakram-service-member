package tiers

import (
	"context"
	"fmt"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type UpdateTierService struct {
	NameTh       string           `json:"name_th"`
	NameEn       string           `json:"name_en"`
	MinSpending  *decimal.Decimal `json:"min_spending"`
	IsActive     *bool            `json:"is_active"`
	DiscountRate *decimal.Decimal `json:"discount_rate"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateTierService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`tiers.svc.update.start`)

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.TierEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", id).Scan(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		if req.NameTh != "" {
			data.NameTh = req.NameTh
		}
		if req.NameEn != "" {
			data.NameEn = req.NameEn
		}
		if req.MinSpending != nil {
			data.MinSpending = *req.MinSpending
		}
		if req.IsActive != nil {
			data.IsActive = *req.IsActive
		}
		if req.DiscountRate != nil {
			data.DiscountRate = *req.DiscountRate
		}

		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_tier",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated tier with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`tiers.svc.update.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_tier",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update tier failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`tiers.svc.update.success`)
	return nil
}
