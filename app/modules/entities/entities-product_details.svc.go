package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.ProductDetailEntity = (*Service)(nil)

func (s *Service) ListProductDetails(ctx context.Context, req *entitiesdto.ListProductDetailsRequest) ([]*ent.ProductDetailEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.ProductDetailEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"description", "material", "dimensions"},
		[]string{"product_id", "description"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetProductDetailByID(ctx context.Context, id uuid.UUID) (*ent.ProductDetailEntity, error) {
	data := new(ent.ProductDetailEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateProductDetail(ctx context.Context, detail *ent.ProductDetailEntity) error {
	_, err := s.db.NewInsert().
		Model(detail).
		Exec(ctx)
	return err
}

func (s *Service) UpdateProductDetail(ctx context.Context, detail *ent.ProductDetailEntity) error {
	_, err := s.db.NewUpdate().
		Model(detail).
		Where("id = ?", detail.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteProductDetail(ctx context.Context, detailID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.ProductDetailEntity{}).
		Where("id = ?", detailID).
		Exec(ctx)
	return err
}
