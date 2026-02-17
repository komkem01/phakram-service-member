package categories

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoCategoryServiceResponses struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	NameTh    string     `json:"name_th"`
	NameEn    string     `json:"name_en"`
	IsActive  bool       `json:"is_active"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoCategoryServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`categories.svc.info.start`)

	data, err := s.db.GetCategoryByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoCategoryServiceResponses{
		ID:        data.ID,
		ParentID:  data.ParentID,
		NameTh:    data.NameTh,
		NameEn:    data.NameEn,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`categories.svc.info.success`)
	return resp, nil
}
