package product_details

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateProductDetailService struct {
	ProductID        uuid.UUID       `json:"product_id"`
	Description      string          `json:"description"`
	Material         string          `json:"material"`
	Dimensions       string          `json:"dimensions"`
	Weight           decimal.Decimal `json:"weight"`
	CareInstructions string          `json:"care_instructions"`
	MemberID         uuid.UUID       `json:"member_id"`
}

func (s *Service) CreateProductDetailService(ctx context.Context, req *CreateProductDetailService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_details.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		productDetail := &ent.ProductDetailEntity{
			ID:               uuid.New(),
			ProductID:        req.ProductID,
			Description:      req.Description,
			Material:         req.Material,
			Dimensions:       req.Dimensions,
			Weight:           req.Weight,
			CareInstructions: req.CareInstructions,
		}
		if _, err := tx.NewInsert().Model(productDetail).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "product_detail",
			ActionID:     &productDetail.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created product detail " + productDetail.ID.String(),
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
	span.AddEvent(`product_details.svc.create.success`)
	return nil
}
