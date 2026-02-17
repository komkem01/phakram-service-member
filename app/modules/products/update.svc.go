package products

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

type UpdateProductService struct {
	CategoryID *string          `json:"category_id"`
	NameTh     string           `json:"name_th"`
	NameEn     string           `json:"name_en"`
	ProductNo  string           `json:"product_no"`
	Price      *decimal.Decimal `json:"price"`
	IsActive   *bool            `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateProductService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`products.svc.update.start`)

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.ProductEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", id).Scan(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		if req.CategoryID != nil && *req.CategoryID != "" {
			categoryID, err := uuid.Parse(*req.CategoryID)
			if err != nil {
				return err
			}
			data.CategoryID = categoryID
		}
		if req.NameTh != "" {
			data.NameTh = req.NameTh
		}
		if req.NameEn != "" {
			data.NameEn = req.NameEn
		}
		if req.ProductNo != "" {
			data.ProductNo = req.ProductNo
		}
		if req.Price != nil {
			data.Price = *req.Price
		}
		if req.IsActive != nil {
			data.IsActive = *req.IsActive
		}

		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_product",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated product with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`products.svc.update.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_product",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update product failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`products.svc.update.success`)
	return nil
}
