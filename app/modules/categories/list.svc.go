package categories

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListCategoryServiceRequest struct {
	base.RequestPaginate
}

type ListCategoryServiceResponses struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	NameTh    string     `json:"name_th"`
	NameEn    string     `json:"name_en"`
	IsActive  bool       `json:"is_active"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListCategoryServiceRequest) ([]*ListCategoryServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`categories.svc.list.start`)

	data, page, err := s.db.ListCategories(ctx, &entitiesdto.ListCategoriesRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListCategoryServiceResponses
	for _, item := range data {
		temp := &ListCategoryServiceResponses{
			ID:        item.ID,
			ParentID:  item.ParentID,
			NameTh:    item.NameTh,
			NameEn:    item.NameEn,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`categories.svc.list.copy`)
	return response, page, nil
}
