package payments

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) InfoService(ctx context.Context, id string) (any, error) {
	paymentID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.db.GetPaymentByID(ctx, paymentID)
}
