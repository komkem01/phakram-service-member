package member_files

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoMemberFileServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoMemberFileServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_files.svc.info.start`)

	data, err := s.db.GetMemberFileByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoMemberFileServiceResponses{
		ID:        data.ID,
		MemberID:  data.MemberID,
		FileID:    data.FileID,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_files.svc.info.success`)
	return resp, nil
}
