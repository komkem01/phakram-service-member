package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.SubDistrictEntity = (*Service)(nil)

func (s *Service) ListSubDistricts(ctx context.Context, req *entitiesdto.ListSubDistrictsRequest) ([]*ent.SubDistrictEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.SubDistrictEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name", "district_id", "is_active"},
		[]string{"created_at", "name", "district_id", "is_active"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetSubDistrictByID(ctx context.Context, id uuid.UUID) (*ent.SubDistrictEntity, error) {
	data := new(ent.SubDistrictEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateSubDistrict(ctx context.Context, subDistrict *ent.SubDistrictEntity) error {
	_, err := s.db.NewInsert().
		Model(subDistrict).
		Exec(ctx)
	return err
}

func (s *Service) UpdateSubDistrict(ctx context.Context, subDistrict *ent.SubDistrictEntity) error {
	_, err := s.db.NewUpdate().
		Model(subDistrict).
		Where("id = ?", subDistrict.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteSubDistrict(ctx context.Context, subDistrictID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.SubDistrictEntity{}).
		Where("id = ?", subDistrictID).
		Exec(ctx)
	return err
}
