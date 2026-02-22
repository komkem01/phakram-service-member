package contact

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type contactMessageReplyRecord struct {
	bun.BaseModel `bun:"table:contact_message_replies"`

	ID               uuid.UUID `bun:"id,pk,type:uuid"`
	ContactMessageID uuid.UUID `bun:"contact_message_id,type:uuid,notnull"`
	SenderRole       string    `bun:"sender_role,notnull"`
	SenderName       string    `bun:"sender_name,notnull"`
	Message          string    `bun:"message,notnull"`
	CreatedAt        time.Time `bun:"created_at,notnull"`
}

type ContactReplyResponse struct {
	ID               string `json:"id"`
	ContactMessageID string `json:"contact_message_id"`
	SenderRole       string `json:"sender_role"`
	SenderName       string `json:"sender_name"`
	Message          string `json:"message"`
	CreatedAt        string `json:"created_at"`
}

type CreateReplyServiceRequest struct {
	ContactMessageID string
	SenderRole       string
	SenderName       string
	Message          string
}

func (s *Service) ListReplies(ctx context.Context, contactMessageID string) ([]*ContactReplyResponse, error) {
	messageID := strings.TrimSpace(contactMessageID)
	messageUUID, err := uuid.Parse(messageID)
	if err != nil {
		return nil, err
	}

	if _, err := s.findContactMessageByID(ctx, messageUUID); err != nil {
		return nil, err
	}

	items := make([]*contactMessageReplyRecord, 0)

	if err := s.bunDB.DB().NewSelect().
		Model(&items).
		Where("contact_message_id = ?", messageUUID).
		OrderExpr("created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	resp := make([]*ContactReplyResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, &ContactReplyResponse{
			ID:               item.ID.String(),
			ContactMessageID: item.ContactMessageID.String(),
			SenderRole:       item.SenderRole,
			SenderName:       item.SenderName,
			Message:          item.Message,
			CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return resp, nil
}

func (s *Service) ListRepliesPublic(ctx context.Context, contactMessageID string, chatToken string) ([]*ContactReplyResponse, error) {
	messageUUID, err := uuid.Parse(strings.TrimSpace(contactMessageID))
	if err != nil {
		return nil, err
	}

	contactMsg, err := s.findContactMessageByIDAndToken(ctx, messageUUID, chatToken)
	if err != nil {
		return nil, err
	}

	items := make([]*contactMessageReplyRecord, 0)

	if err := s.bunDB.DB().NewSelect().
		Model(&items).
		Where("contact_message_id = ?", messageUUID).
		OrderExpr("created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	resp := make([]*ContactReplyResponse, 0, len(items)+1)
	resp = append(resp, &ContactReplyResponse{
		ID:               "origin-" + contactMsg.ID.String(),
		ContactMessageID: contactMsg.ID.String(),
		SenderRole:       "customer",
		SenderName:       contactMsg.Name,
		Message:          contactMsg.Message,
		CreatedAt:        contactMsg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
	for _, item := range items {
		resp = append(resp, &ContactReplyResponse{
			ID:               item.ID.String(),
			ContactMessageID: item.ContactMessageID.String(),
			SenderRole:       item.SenderRole,
			SenderName:       item.SenderName,
			Message:          item.Message,
			CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return resp, nil
}

func (s *Service) CreateReply(ctx context.Context, req *CreateReplyServiceRequest) (*ContactReplyResponse, error) {
	contactMessageID, err := uuid.Parse(strings.TrimSpace(req.ContactMessageID))
	if err != nil {
		return nil, err
	}

	if _, err := s.findContactMessageByID(ctx, contactMessageID); err != nil {
		return nil, err
	}

	now := time.Now()
	record := &contactMessageReplyRecord{
		ID:               uuid.New(),
		ContactMessageID: contactMessageID,
		SenderRole:       strings.TrimSpace(req.SenderRole),
		SenderName:       strings.TrimSpace(req.SenderName),
		Message:          strings.TrimSpace(req.Message),
		CreatedAt:        now,
	}

	if _, err := s.bunDB.DB().NewInsert().Model(record).Exec(ctx); err != nil {
		return nil, err
	}

	return &ContactReplyResponse{
		ID:               record.ID.String(),
		ContactMessageID: record.ContactMessageID.String(),
		SenderRole:       record.SenderRole,
		SenderName:       record.SenderName,
		Message:          record.Message,
		CreatedAt:        record.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) CreateReplyPublic(ctx context.Context, contactMessageID string, chatToken string, message string) (*ContactReplyResponse, error) {
	messageUUID, err := uuid.Parse(strings.TrimSpace(contactMessageID))
	if err != nil {
		return nil, err
	}

	contactMsg, err := s.findContactMessageByIDAndToken(ctx, messageUUID, chatToken)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	record := &contactMessageReplyRecord{
		ID:               uuid.New(),
		ContactMessageID: messageUUID,
		SenderRole:       "customer",
		SenderName:       contactMsg.Name,
		Message:          strings.TrimSpace(message),
		CreatedAt:        now,
	}

	if _, err := s.bunDB.DB().NewInsert().Model(record).Exec(ctx); err != nil {
		return nil, err
	}

	_, _ = s.bunDB.DB().NewUpdate().
		Model((*contactMessageRecord)(nil)).
		Set("is_read = FALSE").
		Set("read_at = NULL").
		Set("updated_at = ?", now).
		Where("id = ?", messageUUID).
		Exec(ctx)

	return &ContactReplyResponse{
		ID:               record.ID.String(),
		ContactMessageID: record.ContactMessageID.String(),
		SenderRole:       record.SenderRole,
		SenderName:       record.SenderName,
		Message:          record.Message,
		CreatedAt:        record.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) findContactMessageByID(ctx context.Context, id uuid.UUID) (*contactMessageRecord, error) {
	message := &contactMessageRecord{}
	if err := s.bunDB.DB().NewSelect().
		Model(message).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *Service) findContactMessageByIDAndToken(ctx context.Context, id uuid.UUID, chatToken string) (*contactMessageRecord, error) {
	token := strings.TrimSpace(chatToken)
	if token == "" {
		return nil, fmt.Errorf("chat token is required")
	}

	message := &contactMessageRecord{}
	if err := s.bunDB.DB().NewSelect().
		Model(message).
		Where("id = ?", id).
		Where("access_token = ?", token).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	return message, nil
}
