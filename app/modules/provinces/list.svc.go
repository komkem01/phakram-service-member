package provinces

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListProvinceServiceRequest struct {
	base.RequestPaginate
}

type ListProvinceServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListProvinceServiceRequest) ([]*ListProvinceServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`provinces.svc.list.start`)

	data, page, err := s.db.ListProvinces(ctx, &entitiesdto.ListProvincesRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListProvinceServiceResponses
	for _, item := range data {
		temp := &ListProvinceServiceResponses{
			ID:        item.ID,
			Name:      item.Name,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`provinces.svc.list.copy`)
	return response, page, nil
}
