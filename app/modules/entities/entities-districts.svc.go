package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.DistrictEntity = (*Service)(nil)

func (s *Service) ListDistricts(ctx context.Context, req *entitiesdto.ListDistrictsRequest) ([]*ent.DistrictEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.DistrictEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name", "province_id", "is_active"},
		[]string{"created_at", "name", "province_id", "is_active"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetDistrictByID(ctx context.Context, id uuid.UUID) (*ent.DistrictEntity, error) {
	data := new(ent.DistrictEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateDistrict(ctx context.Context, district *ent.DistrictEntity) error {
	_, err := s.db.NewInsert().
		Model(district).
		Exec(ctx)
	return err
}

func (s *Service) UpdateDistrict(ctx context.Context, district *ent.DistrictEntity) error {
	_, err := s.db.NewUpdate().
		Model(district).
		Where("id = ?", district.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteDistrict(ctx context.Context, districtID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.DistrictEntity{}).
		Where("id = ?", districtID).
		Exec(ctx)
	return err
}
