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

var _ entitiesinf.ProductEntity = (*Service)(nil)

func (s *Service) ListProducts(ctx context.Context, req *entitiesdto.ListProductsRequest) ([]*ent.ProductEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.ProductEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name_th", "name_en", "product_no", "is_active"},
		[]string{"created_at", "name_th", "name_en", "product_no", "price"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetProductByID(ctx context.Context, id uuid.UUID) (*ent.ProductEntity, error) {
	data := new(ent.ProductEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateProduct(ctx context.Context, product *ent.ProductEntity) error {
	_, err := s.db.NewInsert().
		Model(product).
		Exec(ctx)
	return err
}

func (s *Service) UpdateProduct(ctx context.Context, product *ent.ProductEntity) error {
	_, err := s.db.NewUpdate().
		Model(product).
		Where("id = ?", product.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.ProductEntity{}).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", productID).
		Exec(ctx)
	return err
}
