package products

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateProductService struct {
	CategoryID uuid.UUID       `json:"category_id"`
	NameTh     string          `json:"name_th"`
	NameEn     string          `json:"name_en"`
	ProductNo  string          `json:"product_no"`
	Price      decimal.Decimal `json:"price"`
	IsActive   bool            `json:"is_active"`
	MemberID   uuid.UUID       `json:"member_id"`
}

func (s *Service) CreateProductService(ctx context.Context, req *CreateProductService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`products.svc.create.start`)

	nameForLog := req.NameEn
	if nameForLog == "" {
		nameForLog = req.NameTh
	}

	productNo, err := utils.GenerateProductNo()
	if err != nil {
		return err
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		product := &ent.ProductEntity{
			ID:         uuid.New(),
			CategoryID: req.CategoryID,
			NameTh:     req.NameTh,
			NameEn:     req.NameEn,
			ProductNo:  productNo,
			Price:      req.Price,
			IsActive:   req.IsActive,
		}
		if _, err := tx.NewInsert().Model(product).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "product",
			ActionID:     &product.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created product " + nameForLog,
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
	span.AddEvent(`products.svc.create.success`)
	return nil
}
