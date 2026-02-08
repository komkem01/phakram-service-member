package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.StatusEntity = (*Service)(nil)

func (s *Service) ListStatuses(ctx context.Context, req *entitiesdto.ListStatusesRequest) ([]*ent.StatusEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.StatusEntity, 0)

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

func (s *Service) GetStatusByID(ctx context.Context, id uuid.UUID) (*ent.StatusEntity, error) {
	data := new(ent.StatusEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateStatus(ctx context.Context, status *ent.StatusEntity) error {
	_, err := s.db.NewInsert().
		Model(status).
		Exec(ctx)
	return err
}

func (s *Service) UpdateStatus(ctx context.Context, status *ent.StatusEntity) error {
	_, err := s.db.NewUpdate().
		Model(status).
		Where("id = ?", status.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteStatus(ctx context.Context, statusID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.StatusEntity{}).
		Where("id = ?", statusID).
		Exec(ctx)
	return err
}
