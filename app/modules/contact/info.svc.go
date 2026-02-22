package contact

import (
	"context"
	"strings"
)

func (s *Service) Info(ctx context.Context, id string) (*ListContactServiceResponse, error) {
	message := &contactMessageRecord{}
	if err := s.bunDB.DB().NewSelect().
		Model(message).
		Where("id = ?", strings.TrimSpace(id)).
		Scan(ctx); err != nil {
		return nil, err
	}

	var sentAt *string
	if message.SentAt != nil {
		v := message.SentAt.Format("2006-01-02T15:04:05Z07:00")
		sentAt = &v
	}

	var readAt *string
	if message.ReadAt != nil {
		v := message.ReadAt.Format("2006-01-02T15:04:05Z07:00")
		readAt = &v
	}

	var sendError *string
	if strings.TrimSpace(message.SendError) != "" {
		v := message.SendError
		sendError = &v
	}

	return &ListContactServiceResponse{
		ID:         message.ID.String(),
		Name:       message.Name,
		Email:      message.Email,
		Subject:    message.Subject,
		Message:    message.Message,
		SendStatus: message.SendStatus,
		IsRead:     message.IsRead,
		SendError:  sendError,
		SentAt:     sentAt,
		ReadAt:     readAt,
		CreatedAt:  message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  message.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
