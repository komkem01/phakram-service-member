package storages

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"strings"

	"github.com/google/uuid"
)

func (s *Service) ListService(ctx context.Context, refID uuid.UUID) ([]*ent.StorageEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`storages.svc.list.start`)

	data, err := s.db.ListStoragesByRefID(ctx, refID)
	if err != nil {
		return nil, err
	}

	if s.railwayStorage != nil {
		for _, item := range data {
			if item == nil {
				continue
			}
			resolved, resolveErr := s.railwayStorage.ResolveObjectURL(ctx, item.FilePath)
			if resolveErr == nil && strings.TrimSpace(resolved) != "" {
				item.FilePath = resolved
			}
		}
	}

	span.AddEvent(`storages.svc.list.success`)
	return data, nil
}
