package product_stocks

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoProductStockServiceResponses struct {
	ID          uuid.UUID       `json:"id"`
	ProductID   uuid.UUID       `json:"product_id"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	StockAmount int             `json:"stock_amount"`
	Remaining   int             `json:"remaining"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoProductStockServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_stocks.svc.info.start`)

	data, err := s.db.GetProductStockByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoProductStockServiceResponses{
		ID:          data.ID,
		ProductID:   data.ProductID,
		UnitPrice:   data.UnitPrice,
		StockAmount: data.StockAmount,
		Remaining:   data.Remaining,
		CreatedAt:   data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`product_stocks.svc.info.success`)
	return resp, nil
}
