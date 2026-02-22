package ent

import (
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
