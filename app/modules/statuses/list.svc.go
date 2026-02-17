package statuses

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListStatusServiceRequest struct {
	base.RequestPaginate
}

type ListStatusServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListStatusServiceRequest) ([]*ListStatusServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`statuses.svc.list.start`)

	data, page, err := s.db.ListStatuses(ctx, &entitiesdto.ListStatusesRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListStatusServiceResponses
	for _, item := range data {
		temp := &ListStatusServiceResponses{
			ID:        item.ID,
			NameTh:    item.NameTh,
			NameEn:    item.NameEn,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`statuses.svc.list.copy`)
	return response, page, nil
}
