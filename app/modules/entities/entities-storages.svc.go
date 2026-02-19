package entities

import (
	"context"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StorageEntity = (*Service)(nil)

func (s *Service) UploadStorage(ctx context.Context, storage *ent.StorageEntity) error {
	if storage != nil && storage.UploadedBy != nil && *storage.UploadedBy == uuid.Nil {
		storage.UploadedBy = nil
	}
	_, err := s.db.NewInsert().
		Model(storage).
		Exec(ctx)
	return err
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

func (s *Service) ListStoragesByRefID(ctx context.Context, refID uuid.UUID) ([]*ent.StorageEntity, error) {
	data := make([]*ent.StorageEntity, 0)
	err := s.db.NewSelect().
		Model(&data).
		Where("ref_id = ?", refID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) DeleteStorageByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.StorageEntity{}).
		Set("is_active = false").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (s *Service) DeleteStoragesByRefID(ctx context.Context, refID uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model(&ent.StorageEntity{}).
		Set("is_active = false").
		Where("ref_id = ?", refID).
		Exec(ctx)
	return err
}

func (s *Service) UpdateStatusStorage(ctx context.Context, id uuid.UUID, req *ent.StorageEntity) error {
	storage := new(ent.StorageEntity)
	storage.IsActive = req.IsActive

	_, err := s.db.NewUpdate().
		Model(storage).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
