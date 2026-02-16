package prefixes

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListPrefixServiceRequest struct {
	base.RequestPaginate
}

type ListPrefixServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameEn    string    `json:"name_en"`
	GenderID  uuid.UUID `json:"gender_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListPrefixServiceRequest) ([]*ListPrefixServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`prefixes.svc.list.start`)

	data, page, err := s.db.ListPrefixes(ctx, &entitiesdto.ListPrefixesRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListPrefixServiceResponses
	for _, item := range data {
		temp := &ListPrefixServiceResponses{
			ID:        item.ID,
			NameTh:    item.NameTh,
			NameEn:    item.NameEn,
			GenderID:  item.GenderID,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`prefixes.svc.list.copy`)
	return response, page, nil
}
