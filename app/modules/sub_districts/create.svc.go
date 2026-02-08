package sub_districts

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateSubDistrictService struct {
	DistrictID uuid.UUID `json:"district_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
}

func (s *Service) CreateSubDistrictService(ctx context.Context, req *CreateSubDistrictService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`sub_districts.svc.create.start`)

	subDistrict := &ent.SubDistrictEntity{
		ID:         uuid.New(),
		DistrictID: req.DistrictID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}
	if err := s.db.CreateSubDistrict(ctx, subDistrict); err != nil {
		return err
	}
	span.AddEvent(`sub_districts.svc.create.success`)
	return nil
}
