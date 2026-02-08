package provinces

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateProvinceService struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (s *Service) CreateProvinceService(ctx context.Context, req *CreateProvinceService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`provinces.svc.create.start`)

	province := &ent.ProvinceEntity{
		ID:       uuid.New(),
		Name:     req.Name,
		IsActive: req.IsActive,
	}
	if err := s.db.CreateProvince(ctx, province); err != nil {
		return err
	}
	span.AddEvent(`provinces.svc.create.success`)
	return nil
}
