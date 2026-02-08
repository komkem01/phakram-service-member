package member_transactions

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type CreateMemberTransactionService struct {
	MemberID uuid.UUID `json:"member_id"`
	Action   string    `json:"action"`
	Details  string    `json:"details"`
}

func (s *Service) CreateMemberTransactionService(ctx context.Context, req *CreateMemberTransactionService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_transactions.svc.create.start`)

	transaction := &ent.MemberTransactionEntity{
		ID:       uuid.New(),
		MemberID: req.MemberID,
		Action:   ent.ActionTypeEnum(req.Action),
		Details:  req.Details,
	}
	if err := s.db.CreateMemberTransaction(ctx, transaction); err != nil {
		return err
	}
	span.AddEvent(`member_transactions.svc.create.success`)
	return nil
}
