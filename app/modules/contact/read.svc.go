package contact

import (
	"context"
	"time"
)

func (s *Service) MarkRead(ctx context.Context, id string, isRead bool) error {
	query := s.bunDB.DB().NewUpdate().
		Model((*contactMessageRecord)(nil)).
		Set("is_read = ?", isRead).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id)

	if isRead {
		now := time.Now()
		query = query.Set("read_at = ?", now)
	} else {
		query = query.Set("read_at = NULL")
	}

	_, err := query.Exec(ctx)
	return err
}
