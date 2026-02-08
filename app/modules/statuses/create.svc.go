package statuses

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateStatusService struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
}

func (s *Service) CreateStatusService(ctx context.Context, req *CreateStatusService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`statuses.svc.create.start`)

	status := &ent.StatusEntity{
		ID:     uuid.New(),
		NameTh: req.NameTh,
		NameEn: req.NameEn,
	}
	if err := s.db.CreateStatus(ctx, status); err != nil {
		return err
	}
	span.AddEvent(`statuses.svc.create.success`)
	return nil
}
