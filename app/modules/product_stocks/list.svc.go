package product_stocks

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListProductStockServiceRequest struct {
	base.RequestPaginate
}

type ListProductStockServiceResponses struct {
	ID          uuid.UUID       `json:"id"`
	ProductID   uuid.UUID       `json:"product_id"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	StockAmount int             `json:"stock_amount"`
	Remaining   int             `json:"remaining"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListProductStockServiceRequest) ([]*ListProductStockServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_stocks.svc.list.start`)

	data, page, err := s.db.ListProductStocks(ctx, &entitiesdto.ListProductStocksRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListProductStockServiceResponses
	for _, item := range data {
		temp := &ListProductStockServiceResponses{
			ID:          item.ID,
			ProductID:   item.ProductID,
			UnitPrice:   item.UnitPrice,
			StockAmount: item.StockAmount,
			Remaining:   item.Remaining,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`product_stocks.svc.list.copy`)
	return response, page, nil
}
