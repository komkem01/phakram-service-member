package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.GenderEntity = (*Service)(nil)

func (s *Service) ListGenders(ctx context.Context, req *entitiesdto.ListGendersRequest) ([]*ent.GenderEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.GenderEntity, 0)

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

func (s *Service) GetGenderByID(ctx context.Context, id uuid.UUID) (*ent.GenderEntity, error) {
	data := new(ent.GenderEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateGender(ctx context.Context, gender *ent.GenderEntity) error {
	_, err := s.db.NewInsert().
		Model(gender).
		Exec(ctx)
	return err
}

func (s *Service) UpdateGender(ctx context.Context, gender *ent.GenderEntity) error {
	_, err := s.db.NewUpdate().
		Model(gender).
		Where("id = ?", gender.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteGender(ctx context.Context, genderID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.GenderEntity{}).
		Where("id = ?", genderID).
		Exec(ctx)
	return err
}
