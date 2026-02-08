package districts

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateDistrictService struct {
	ProvinceID uuid.UUID `json:"province_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
}

func (s *Service) CreateDistrictService(ctx context.Context, req *CreateDistrictService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`districts.svc.create.start`)

	district := &ent.DistrictEntity{
		ID:         uuid.New(),
		ProvinceID: req.ProvinceID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}
	if err := s.db.CreateDistrict(ctx, district); err != nil {
		return err
	}
	span.AddEvent(`districts.svc.create.success`)
	return nil
}
