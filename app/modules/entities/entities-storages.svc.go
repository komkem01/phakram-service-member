package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.StorageEntity = (*Service)(nil)

func (s *Service) ListStorages(ctx context.Context, req *entitiesdto.ListStoragesRequest) ([]*ent.StorageEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.StorageEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"file_name", "file_path", "file_type", "related_entity"},
		[]string{"created_at", "file_name", "file_type", "related_entity"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetStorageByID(ctx context.Context, id uuid.UUID) (*ent.StorageEntity, error) {
	data := new(ent.StorageEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateStorage(ctx context.Context, storage *ent.StorageEntity) error {
	_, err := s.db.NewInsert().
		Model(storage).
		Exec(ctx)
	return err
}

func (s *Service) UpdateStorage(ctx context.Context, storage *ent.StorageEntity) error {
	_, err := s.db.NewUpdate().
		Model(storage).
		Where("id = ?", storage.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteStorage(ctx context.Context, storageID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.StorageEntity{}).
		Where("id = ?", storageID).
		Exec(ctx)
	return err
}
