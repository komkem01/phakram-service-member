package contact

import "context"

func (s *Service) CountUnread(ctx context.Context) (int64, error) {
	count, err := s.bunDB.DB().NewSelect().
		Model((*contactMessageRecord)(nil)).
		Where("is_read = ?", false).
		Count(ctx)
	if err != nil {
		return 0, err
	}

	return int64(count), nil
}
