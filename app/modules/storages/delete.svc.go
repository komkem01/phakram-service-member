package storages

import (
	"context"
	"phakram/app/utils"

	"github.com/google/uuid"
)

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.delete.start`)

	if err := s.db.DeleteStorageByID(ctx, id); err != nil {
		return err
	}

	span.AddEvent(`storages.svc.delete.success`)
	return nil
}
