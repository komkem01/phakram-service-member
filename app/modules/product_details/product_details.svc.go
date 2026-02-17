package productdetails

import (
	"context"
	"phakram/app/modules/entities/ent"

	"github.com/google/uuid"
)

func (s *Service) GetByProductID(ctx context.Context, productID uuid.UUID) (*ent.ProductDetailEntity, error) {
	data := new(ent.ProductDetailEntity)
	err := s.bunDB.DB().NewSelect().Model(data).Where("product_id = ?", productID).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateByProductID(ctx context.Context, productID uuid.UUID, payload *ent.ProductDetailEntity) error {
	payload.ID = uuid.New()
	payload.ProductID = productID
	_, err := s.bunDB.DB().NewInsert().Model(payload).Exec(ctx)
	return err
}

func (s *Service) UpdateByProductID(ctx context.Context, productID uuid.UUID, payload *ent.ProductDetailEntity) error {
	current, err := s.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}
	current.Description = payload.Description
	current.Material = payload.Material
	current.Dimensions = payload.Dimensions
	current.Weight = payload.Weight
	current.CareInstructions = payload.CareInstructions
	_, err = s.bunDB.DB().NewUpdate().Model(current).Where("id = ?", current.ID).Exec(ctx)
	return err
}

func (s *Service) DeleteByProductID(ctx context.Context, productID uuid.UUID) error {
	_, err := s.bunDB.DB().NewDelete().Model((*ent.ProductDetailEntity)(nil)).Where("product_id = ?", productID).Exec(ctx)
	return err
}
