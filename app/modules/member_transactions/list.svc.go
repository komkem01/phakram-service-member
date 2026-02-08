package member_transactions

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberTransactionServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListMemberTransactionServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	CreatedAt string    `json:"created_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListMemberTransactionServiceRequest) ([]*ListMemberTransactionServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_transactions.svc.list.start`)

	data, page, err := s.db.ListMemberTransactions(ctx, &entitiesdto.ListMemberTransactionsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListMemberTransactionServiceResponses
	for _, item := range data {
		temp := &ListMemberTransactionServiceResponses{
			ID:        item.ID,
			MemberID:  item.MemberID,
			Action:    string(item.Action),
			Details:   item.Details,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`member_transactions.svc.list.copy`)
	return response, page, nil
}
