package storages

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

func (s *Service) ListService(ctx context.Context, refID uuid.UUID) ([]*ent.StorageEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.list.start`)

	data, err := s.db.ListStoragesByRefID(ctx, refID)
	if err != nil {
		return nil, err
	}

	span.AddEvent(`storages.svc.list.success`)
	return data, nil
}
