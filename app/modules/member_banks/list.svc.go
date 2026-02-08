package member_banks

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberBankServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListMemberBankServiceResponses struct {
	ID          uuid.UUID `json:"id"`
	MemberID    uuid.UUID `json:"member_id"`
	BankID      uuid.UUID `json:"bank_id"`
	BankNo      string    `json:"bank_no"`
	FirstnameTh string    `json:"firstname_th"`
	LastnameTh  string    `json:"lastname_th"`
	FirstnameEn string    `json:"firstname_en"`
	LastnameEn  string    `json:"lastname_en"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListMemberBankServiceRequest) ([]*ListMemberBankServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_banks.svc.list.start`)

	data, page, err := s.db.ListMemberBanks(ctx, &entitiesdto.ListMemberBanksRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListMemberBankServiceResponses
	for _, item := range data {
		temp := &ListMemberBankServiceResponses{
			ID:          item.ID,
			MemberID:    item.MemberID,
			BankID:      item.BankID,
			BankNo:      item.BankNo,
			FirstnameTh: item.FirstnameTh,
			LastnameTh:  item.LastnameTh,
			FirstnameEn: item.FirstnameEn,
			LastnameEn:  item.LastnameEn,
			IsSystem:    item.IsSystem,
			IsActive:    item.IsActive,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`member_banks.svc.list.copy`)
	return response, page, nil
}
