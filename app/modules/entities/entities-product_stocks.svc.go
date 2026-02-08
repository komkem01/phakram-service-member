package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.ProductStockEntity = (*Service)(nil)

func (s *Service) ListProductStocks(ctx context.Context, req *entitiesdto.ListProductStocksRequest) ([]*ent.ProductStockEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.ProductStockEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"product_id", "stock_amount", "remaining"},
		[]string{"created_at", "product_id", "stock_amount", "remaining"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetProductStockByID(ctx context.Context, id uuid.UUID) (*ent.ProductStockEntity, error) {
	data := new(ent.ProductStockEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateProductStock(ctx context.Context, stock *ent.ProductStockEntity) error {
	_, err := s.db.NewInsert().
		Model(stock).
		Exec(ctx)
	return err
}

func (s *Service) UpdateProductStock(ctx context.Context, stock *ent.ProductStockEntity) error {
	_, err := s.db.NewUpdate().
		Model(stock).
		Where("id = ?", stock.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteProductStock(ctx context.Context, stockID uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.ProductStockEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", stockID).
		Exec(ctx)
	return err
}
