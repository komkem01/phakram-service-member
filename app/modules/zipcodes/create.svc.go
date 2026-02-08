package zipcodes

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateZipcodeService struct {
	SubDistrictsID uuid.UUID `json:"sub_districts_id"`
	Name           string    `json:"name"`
	IsActive       bool      `json:"is_active"`
}

func (s *Service) CreateZipcodeService(ctx context.Context, req *CreateZipcodeService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`zipcodes.svc.create.start`)

	zipcode := &ent.ZipcodeEntity{
		ID:             uuid.New(),
		SubDistrictsID: req.SubDistrictsID,
		Name:           req.Name,
		IsActive:       req.IsActive,
	}
	if err := s.db.CreateZipcode(ctx, zipcode); err != nil {
		return err
	}
	span.AddEvent(`zipcodes.svc.create.success`)
	return nil
}
