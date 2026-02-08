package districts

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListDistrictServiceRequest struct {
	base.RequestPaginate
}

type ListDistrictServiceResponses struct {
	ID         uuid.UUID `json:"id"`
	ProvinceID uuid.UUID `json:"province_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListDistrictServiceRequest) ([]*ListDistrictServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`districts.svc.list.start`)

	data, page, err := s.db.ListDistricts(ctx, &entitiesdto.ListDistrictsRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListDistrictServiceResponses
	for _, item := range data {
		temp := &ListDistrictServiceResponses{
			ID:         item.ID,
			ProvinceID: item.ProvinceID,
			Name:       item.Name,
			IsActive:   item.IsActive,
			CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`districts.svc.list.copy`)
	return response, page, nil
}
