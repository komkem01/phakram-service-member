package products

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

type CreateProductService struct {
	CategoryID string          `json:"category_id"`
	NameTh     string          `json:"name_th"`
	NameEn     string          `json:"name_en"`
	Price      decimal.Decimal `json:"price"`
	IsActive   *bool           `json:"is_active"`
}

func (s *Service) CreateProductService(ctx context.Context, req *CreateProductService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`products.svc.create.start`)

	id := uuid.New()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return err
	}

	productNo, err := s.generateUniqueProductNo(ctx)
	if err != nil {
		return err
	}

	product := &ent.ProductEntity{
		ID:         id,
		CategoryID: categoryID,
		NameTh:     req.NameTh,
		NameEn:     req.NameEn,
		ProductNo:  productNo,
		Price:      req.Price,
		IsActive:   isActive,
	}
	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(product).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_product",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created product with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`products.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_product",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create product failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`products.svc.create.success`)
	return nil
}

func (s *Service) generateUniqueProductNo(ctx context.Context) (string, error) {
	for range 20 {
		code, err := utils.GenerateProductNo()
		if err != nil {
			return "", err
		}

		exists, err := s.bunDB.DB().NewSelect().
			Model((*ent.ProductEntity)(nil)).
			Where("product_no = ?", code).
			Where("deleted_at IS NULL").
			Exists(ctx)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", fmt.Errorf("unable to generate unique product number")
}
