package orders

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListOrderServiceRequest struct {
	base.RequestPaginate
	MemberID  uuid.UUID
	Search    string
	Status    string
	StartDate int64
	EndDate   int64
}

type ListOrderServiceResponses struct {
	ID             uuid.UUID       `json:"id"`
	OrderNo        string          `json:"order_no"`
	MemberID       uuid.UUID       `json:"member_id"`
	PaymentID      uuid.UUID       `json:"payment_id"`
	AddressID      uuid.UUID       `json:"address_id"`
	Status         string          `json:"status"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	NetAmount      decimal.Decimal `json:"net_amount"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListOrderServiceRequest) ([]*ListOrderServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.list.start`)

	data, page, err := s.db.ListOrders(ctx, &entitiesdto.ListOrdersRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
		Search:          req.Search,
		Status:          req.Status,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListOrderServiceResponses
	for _, item := range data {
		temp := &ListOrderServiceResponses{
			ID:             item.ID,
			OrderNo:        item.OrderNo,
			MemberID:       item.MemberID,
			PaymentID:      item.PaymentID,
			AddressID:      item.AddressID,
			Status:         string(item.Status),
			TotalAmount:    item.TotalAmount,
			DiscountAmount: item.DiscountAmount,
			NetAmount:      item.NetAmount,
			CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`orders.svc.list.copy`)
	return response, page, nil
}
