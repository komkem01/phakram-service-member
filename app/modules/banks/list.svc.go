package banks

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListBankServiceRequest struct {
	base.RequestPaginate
}

type ListBankServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	NameTh    string    `json:"name_th"`
	NameAbbTh string    `json:"name_abb_th"`
	NameEn    string    `json:"name_en"`
	NameAbbEn string    `json:"name_abb_en"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListBankServiceRequest) ([]*ListBankServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`banks.svc.list.start`)

	data, page, err := s.db.ListBanks(ctx, &entitiesdto.ListBanksRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListBankServiceResponses
	for _, item := range data {
		temp := &ListBankServiceResponses{
			ID:        item.ID,
			NameTh:    item.NameTh,
			NameAbbTh: item.NameAbbTh,
			NameEn:    item.NameEn,
			NameAbbEn: item.NameAbbEn,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`banks.svc.list.copy`)
	return response, page, nil
}
