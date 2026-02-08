package districts

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`districts.svc.delete.start`)

	if err := s.db.DeleteDistrict(ctx, id); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`districts.svc.delete.success`)
	return nil
}
