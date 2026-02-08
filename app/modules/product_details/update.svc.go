package product_details

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type UpdateProductDetailService struct {
	ProductID        uuid.UUID        `json:"product_id"`
	Description      string           `json:"description"`
	Material         string           `json:"material"`
	Dimensions       string           `json:"dimensions"`
	Weight           *decimal.Decimal `json:"weight"`
	CareInstructions string           `json:"care_instructions"`
	MemberID         uuid.UUID        `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateProductDetailService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_details.svc.update.start`)

	data, err := s.db.GetProductDetailByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.ProductID != uuid.Nil {
		data.ProductID = req.ProductID
	}
	if req.Description != "" {
		data.Description = req.Description
	}
	if req.Material != "" {
		data.Material = req.Material
	}
	if req.Dimensions != "" {
		data.Dimensions = req.Dimensions
	}
	if req.Weight != nil {
		data.Weight = *req.Weight
	}
	if req.CareInstructions != "" {
		data.CareInstructions = req.CareInstructions
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
			ActionType:   "product_detail",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated product detail " + data.ID.String(),
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

	span.AddEvent(`product_details.svc.update.success`)
	return nil
}
