package genders

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateGenderService struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
	IsActive bool `json:"is_active"`
}

func (s *Service) CreateGenderService(ctx context.Context, req *CreateGenderService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`genders.svc.create.start`)

	id := uuid.New()

	// Create gender
	gender := &ent.GenderEntity{
		ID:     id,
		NameTh: req.NameTh,
		NameEn: req.NameEn,
		IsActive: req.IsActive,
	}
	if err := s.db.CreateGender(ctx, gender); err != nil {
		return err
	}
	span.AddEvent(`genders.svc.create.gender_created`)

	span.AddEvent(`genders.svc.create.success`)
	return nil
}
