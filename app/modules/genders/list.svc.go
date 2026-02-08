package genders

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListGenderServiceRequest struct {
	base.RequestPaginate
}

type ListGenderServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	CreatedAt string    `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

func (s *Service) ListService(ctx context.Context, req *ListGenderServiceRequest) ([]*ListGenderServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`genders.svc.list.start`)

	data, page, err := s.db.ListGenders(ctx, &entitiesdto.ListGendersRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListGenderServiceResponses
	for _, item := range data {

		temp := &ListGenderServiceResponses{
			ID:        item.ID,
			NameTh:    item.NameTh,
			NameEn:    item.NameEn,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			IsActive:  item.IsActive,
		}
		response = append(response, temp)
	}
	span.AddEvent(`genders.svc.list.copy`)
	return response, page, nil
}
