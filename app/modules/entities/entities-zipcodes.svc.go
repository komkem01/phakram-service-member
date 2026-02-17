package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.ZipcodeEntity = (*Service)(nil)

func (s *Service) ListZipcodes(ctx context.Context, req *entitiesdto.ListZipcodesRequest) ([]*ent.ZipcodeEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.ZipcodeEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name", "sub_district_id", "is_active"},
		[]string{"created_at", "name", "sub_district_id", "is_active"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetZipcodeByID(ctx context.Context, id uuid.UUID) (*ent.ZipcodeEntity, error) {
	data := new(ent.ZipcodeEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateZipcode(ctx context.Context, zipcode *ent.ZipcodeEntity) error {
	_, err := s.db.NewInsert().
		Model(zipcode).
		Exec(ctx)
	return err
}

func (s *Service) UpdateZipcode(ctx context.Context, zipcode *ent.ZipcodeEntity) error {
	_, err := s.db.NewUpdate().
		Model(zipcode).
		Where("id = ?", zipcode.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteZipcode(ctx context.Context, zipcodeID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.ZipcodeEntity{}).
		Where("id = ?", zipcodeID).
		Exec(ctx)
	return err
}
