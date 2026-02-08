package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.ProductFileEntity = (*Service)(nil)

func (s *Service) ListProductFiles(ctx context.Context, req *entitiesdto.ListProductFilesRequest) ([]*ent.ProductFileEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.ProductFileEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"product_id", "file_id"},
		[]string{"created_at", "product_id", "file_id"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetProductFileByID(ctx context.Context, id uuid.UUID) (*ent.ProductFileEntity, error) {
	data := new(ent.ProductFileEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetProductFileByProductID(ctx context.Context, productID uuid.UUID) (*ent.ProductFileEntity, error) {
	data := new(ent.ProductFileEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateProductFile(ctx context.Context, file *ent.ProductFileEntity) error {
	_, err := s.db.NewInsert().
		Model(file).
		Exec(ctx)
	return err
}

func (s *Service) UpdateProductFile(ctx context.Context, file *ent.ProductFileEntity) error {
	_, err := s.db.NewUpdate().
		Model(file).
		Where("id = ?", file.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteProductFile(ctx context.Context, fileID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.ProductFileEntity{}).
		Where("id = ?", fileID).
		Exec(ctx)
	return err
}
