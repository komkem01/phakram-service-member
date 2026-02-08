package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.CategoryEntity = (*Service)(nil)

func (s *Service) ListCategories(ctx context.Context, req *entitiesdto.ListCategoriesRequest) ([]*ent.CategoryEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.CategoryEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name_th", "name_en"},
		[]string{"created_at", "name_th", "name_en"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, id uuid.UUID) (*ent.CategoryEntity, error) {
	data := new(ent.CategoryEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateCategory(ctx context.Context, category *ent.CategoryEntity) error {
	_, err := s.db.NewInsert().
		Model(category).
		Exec(ctx)
	return err
}

func (s *Service) UpdateCategory(ctx context.Context, category *ent.CategoryEntity) error {
	_, err := s.db.NewUpdate().
		Model(category).
		Where("id = ?", category.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.CategoryEntity{}).
		Where("id = ?", categoryID).
		Exec(ctx)
	return err
}
