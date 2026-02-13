package storages

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*ent.StorageEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.info.start`)

	data, err := s.db.GetStorageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	span.AddEvent(`storages.svc.info.success`)
	return data, nil
}
