package carts

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListCartServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListCartServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListCartServiceRequest) ([]*ListCartServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.list.start`)

	data, page, err := s.db.ListCarts(ctx, &entitiesdto.ListCartsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListCartServiceResponses
	for _, item := range data {
		temp := &ListCartServiceResponses{
			ID:        item.ID,
			MemberID:  item.MemberID,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`carts.svc.list.copy`)
	return response, page, nil
}
