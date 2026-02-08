package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.PrefixEntity = (*Service)(nil)

func (s *Service) ListPrefixes(ctx context.Context, req *entitiesdto.ListPrefixesRequest) ([]*ent.PrefixEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.PrefixEntity, 0)

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

func (s *Service) GetPrefixByID(ctx context.Context, id uuid.UUID) (*ent.PrefixEntity, error) {
	data := new(ent.PrefixEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreatePrefix(ctx context.Context, prefix *ent.PrefixEntity) error {
	_, err := s.db.NewInsert().
		Model(prefix).
		Exec(ctx)
	return err
}

func (s *Service) UpdatePrefix(ctx context.Context, prefix *ent.PrefixEntity) error {
	_, err := s.db.NewUpdate().
		Model(prefix).
		Where("id = ?", prefix.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeletePrefix(ctx context.Context, prefixID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.PrefixEntity{}).
		Where("id = ?", prefixID).
		Exec(ctx)
	return err
}
