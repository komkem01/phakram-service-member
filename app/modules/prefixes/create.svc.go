package prefixes

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreatePrefixService struct {
	NameTh   string    `json:"name_th"`
	NameEn   string    `json:"name_en"`
	GenderID uuid.UUID `json:"gender_id"`
	IsActive bool      `json:"is_active"`
}

func (s *Service) CreatePrefixService(ctx context.Context, req *CreatePrefixService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`prefixes.svc.create.start`)

	id := uuid.New()

	prefix := &ent.PrefixEntity{
		ID:       id,
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		GenderID: req.GenderID,
		IsActive: req.IsActive,
	}
	if err := s.db.CreatePrefix(ctx, prefix); err != nil {
		return err
	}
	span.AddEvent(`prefixes.svc.create.prefix_created`)

	span.AddEvent(`prefixes.svc.create.success`)
	return nil
}
