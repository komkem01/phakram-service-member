package entities

import (
	"context"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"time"
)

var _ entitiesinf.MemberTransactionEntity = (*Service)(nil)

func (s *Service) CreateMemberTransaction(ctx context.Context, memberTransaction *ent.MemberTransactionEntity) error {
	data := ent.MemberTransactionEntity{
		ID:        memberTransaction.ID,
		MemberID:  memberTransaction.MemberID,
		Action:    memberTransaction.Action,
		Details:   memberTransaction.Details,
		CreatedAt: time.Now(),
	}
	_, err := s.db.NewInsert().Model(&data).Exec(ctx)
	return err
}
