package products

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

type UpdateProductService struct {
	CategoryID uuid.UUID        `json:"category_id"`
	NameTh     string           `json:"name_th"`
	NameEn     string           `json:"name_en"`
	ProductNo  string           `json:"product_no"`
	Price      *decimal.Decimal `json:"price"`
	IsActive   *bool            `json:"is_active"`
	MemberID   uuid.UUID        `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateProductService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`products.svc.update.start`)

	data, err := s.db.GetProductByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.CategoryID != uuid.Nil {
		data.CategoryID = req.CategoryID
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

	nameForLog := data.NameEn
	if nameForLog == "" {
		nameForLog = data.NameTh
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
			ActionType:   "product",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated product " + nameForLog,
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

	span.AddEvent(`products.svc.update.success`)
	return nil
}
