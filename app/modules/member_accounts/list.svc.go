package member_accounts

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberAccountServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListMemberAccountServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListMemberAccountServiceRequest) ([]*ListMemberAccountServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_accounts.svc.list.start`)

	data, page, err := s.db.ListMemberAccounts(ctx, &entitiesdto.ListMemberAccountsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListMemberAccountServiceResponses
	for _, item := range data {
		temp := &ListMemberAccountServiceResponses{
			ID:        item.ID,
			MemberID:  item.MemberID,
			Email:     item.Email,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`member_accounts.svc.list.copy`)
	return response, page, nil
}
