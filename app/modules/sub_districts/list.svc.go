package sub_districts

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListSubDistrictServiceRequest struct {
	base.RequestPaginate
}

type ListSubDistrictServiceResponses struct {
	ID         uuid.UUID `json:"id"`
	DistrictID uuid.UUID `json:"district_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListSubDistrictServiceRequest) ([]*ListSubDistrictServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`sub_districts.svc.list.start`)

	data, page, err := s.db.ListSubDistricts(ctx, &entitiesdto.ListSubDistrictsRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListSubDistrictServiceResponses
	for _, item := range data {
		temp := &ListSubDistrictServiceResponses{
			ID:         item.ID,
			DistrictID: item.DistrictID,
			Name:       item.Name,
			IsActive:   item.IsActive,
			CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`sub_districts.svc.list.copy`)
	return response, page, nil
}
