package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type contactMessageRecord struct {
	bun.BaseModel `bun:"table:contact_messages"`

	ID         uuid.UUID  `bun:"id,pk,type:uuid"`
	Name       string     `bun:"name,notnull"`
	Email      string     `bun:"email,notnull"`
	Subject    string     `bun:"subject,notnull"`
	Message    string     `bun:"message,notnull"`
	SendStatus string     `bun:"send_status,notnull"`
	SendError  string     `bun:"send_error"`
	SentAt     *time.Time `bun:"sent_at"`
	CreatedAt  time.Time  `bun:"created_at,notnull"`
	UpdatedAt  time.Time  `bun:"updated_at,notnull"`
}
