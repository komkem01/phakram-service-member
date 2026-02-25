package storages

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"strings"

	"github.com/google/uuid"
)

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*ent.StorageEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.info.start`)

	data, err := s.db.GetStorageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.railwayStorage != nil {
		resolved, resolveErr := s.railwayStorage.ResolveObjectURL(ctx, data.FilePath)
		if resolveErr == nil && strings.TrimSpace(resolved) != "" {
			data.FilePath = resolved
		}
	}

	span.AddEvent(`storages.svc.info.success`)
	return data, nil
}
