package zipcodes

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListZipcodeServiceRequest struct {
	base.RequestPaginate
}

type ListZipcodeServiceResponses struct {
	ID             uuid.UUID `json:"id"`
	SubDistrictsID uuid.UUID `json:"sub_districts_id"`
	Name           string    `json:"name"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListZipcodeServiceRequest) ([]*ListZipcodeServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`zipcodes.svc.list.start`)

	data, page, err := s.db.ListZipcodes(ctx, &entitiesdto.ListZipcodesRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListZipcodeServiceResponses
	for _, item := range data {
		temp := &ListZipcodeServiceResponses{
			ID:             item.ID,
			SubDistrictsID: item.SubDistrictsID,
			Name:           item.Name,
			IsActive:       item.IsActive,
			CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`zipcodes.svc.list.copy`)
	return response, page, nil
}
