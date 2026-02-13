package storages

import (
	"context"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, isActive bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.update.start`)

	if err := s.db.UpdateStatusStorage(ctx, id, &ent.StorageEntity{
		IsActive:  isActive,
		UpdatedAt: time.Now(),
	}); err != nil {
		return err
	}

	span.AddEvent(`storages.svc.update.success`)
	return nil
}
