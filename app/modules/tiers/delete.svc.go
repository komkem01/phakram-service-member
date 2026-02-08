package tiers

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`tiers.svc.delete.start`)

	if err := s.db.DeleteTier(ctx, id); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`tiers.svc.delete.success`)
	return nil
}
