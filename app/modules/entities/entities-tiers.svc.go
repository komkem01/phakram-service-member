package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.TierEntity = (*Service)(nil)

func (s *Service) ListTiers(ctx context.Context, req *entitiesdto.ListTiersRequest) ([]*ent.TierEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.TierEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"name_th", "name_en", "is_active"},
		[]string{"created_at", "name_th", "name_en", "is_active"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetTierByID(ctx context.Context, id uuid.UUID) (*ent.TierEntity, error) {
	data := new(ent.TierEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateTier(ctx context.Context, tier *ent.TierEntity) error {
	_, err := s.db.NewInsert().
		Model(tier).
		Exec(ctx)
	return err
}

func (s *Service) UpdateTier(ctx context.Context, tier *ent.TierEntity) error {
	_, err := s.db.NewUpdate().
		Model(tier).
		Where("id = ?", tier.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteTier(ctx context.Context, tierID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.TierEntity{}).
		Where("id = ?", tierID).
		Exec(ctx)
	return err
}
