package productstocks

import (
	"context"
	"phakram/app/modules/entities/ent"
	"time"

	"github.com/google/uuid"
)

func (s *Service) GetByProductID(ctx context.Context, productID uuid.UUID) (*ent.ProductStockEntity, error) {
	data := new(ent.ProductStockEntity)
	err := s.bunDB.DB().NewSelect().Model(data).Where("product_id = ?", productID).Where("deleted_at IS NULL").Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateByProductID(ctx context.Context, productID uuid.UUID, payload *ent.ProductStockEntity) error {
	product := new(ent.ProductEntity)
	if err := s.bunDB.DB().NewSelect().Model(product).Where("id = ?", productID).Where("deleted_at IS NULL").Limit(1).Scan(ctx); err != nil {
		return err
	}
	payload.UnitPrice = product.Price
	payload.ID = uuid.New()
	payload.ProductID = productID
	_, err := s.bunDB.DB().NewInsert().Model(payload).Exec(ctx)
	return err
}

func (s *Service) UpdateByProductID(ctx context.Context, productID uuid.UUID, payload *ent.ProductStockEntity) error {
	current, err := s.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}
	product := new(ent.ProductEntity)
	if err := s.bunDB.DB().NewSelect().Model(product).Where("id = ?", productID).Where("deleted_at IS NULL").Limit(1).Scan(ctx); err != nil {
		return err
	}
	current.UnitPrice = product.Price
	current.StockAmount = payload.StockAmount
	current.Remaining = payload.Remaining
	_, err = s.bunDB.DB().NewUpdate().Model(current).Where("id = ?", current.ID).Exec(ctx)
	return err
}

func (s *Service) DeleteByProductID(ctx context.Context, productID uuid.UUID) error {
	_, err := s.bunDB.DB().NewUpdate().Model((*ent.ProductStockEntity)(nil)).Set("deleted_at = ?", time.Now()).Where("product_id = ?", productID).Where("deleted_at IS NULL").Exec(ctx)
	return err
}
