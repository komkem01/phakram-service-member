package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.ProvinceEntity = (*Service)(nil)

func (s *Service) ListProvinces(ctx context.Context, req *entitiesdto.ListProvincesRequest) ([]*ent.ProvinceEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.ProvinceEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name", "is_active"},
		[]string{"created_at", "name", "is_active"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetProvinceByID(ctx context.Context, id uuid.UUID) (*ent.ProvinceEntity, error) {
	data := new(ent.ProvinceEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateProvince(ctx context.Context, province *ent.ProvinceEntity) error {
	_, err := s.db.NewInsert().
		Model(province).
		Exec(ctx)
	return err
}

func (s *Service) UpdateProvince(ctx context.Context, province *ent.ProvinceEntity) error {
	_, err := s.db.NewUpdate().
		Model(province).
		Where("id = ?", province.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteProvince(ctx context.Context, provinceID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.ProvinceEntity{}).
		Where("id = ?", provinceID).
		Exec(ctx)
	return err
}
