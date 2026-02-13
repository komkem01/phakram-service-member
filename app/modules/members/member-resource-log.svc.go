package members

import (
	"context"
	"fmt"
	"time"

	"phakram/app/modules/entities/ent"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (s *Service) logMemberActionTx(
	ctx context.Context,
	tx bun.Tx,
	memberID uuid.UUID,
	action ent.MemberActionEnum,
	auditAction ent.AuditActionEnum,
	actionType string,
	actionID uuid.UUID,
	actionBy *uuid.UUID,
	detail string,
	now time.Time,
) error {
	memberTx := &ent.MemberTransactionEntity{
		ID:        uuid.New(),
		MemberID:  memberID,
		Action:    action,
		Details:   detail,
		CreatedAt: now,
	}
	if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
		return err
	}

	auditLog := &ent.AuditLogEntity{
		ID:           uuid.New(),
		Action:       auditAction,
		ActionType:   actionType,
		ActionID:     actionID,
		ActionBy:     actionBy,
		Status:       ent.StatusAuditSuccesses,
		ActionDetail: detail,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) logMemberActionFailed(
	ctx context.Context,
	auditAction ent.AuditActionEnum,
	actionType string,
	actionID uuid.UUID,
	actionBy *uuid.UUID,
	now time.Time,
	err error,
) {
	failLog := &ent.AuditLogEntity{
		ID:           uuid.New(),
		Action:       auditAction,
		ActionType:   actionType,
		ActionID:     actionID,
		ActionBy:     actionBy,
		Status:       ent.StatusAuditFailed,
		ActionDetail: fmt.Sprintf("%s failed: %v", actionType, err),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
}
