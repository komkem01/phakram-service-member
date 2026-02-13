package members

import (
	"context"
	"fmt"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.delete.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Table("members").
			Set("deleted_at = ?", now).
			Set("updated_at = ?", now).
			Where("id = ?", id).
			Where("deleted_at IS NULL").
			Exec(ctx)
		if err != nil {
			return err
		}

		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("no rows in result set")
		}

		memberTx := &ent.MemberTransactionEntity{
			ID:        uuid.New(),
			MemberID:  id,
			Action:    ent.MemberActionDeleted,
			Details:   "member deleted",
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionDeleted,
			ActionType:   "delete_member",
			ActionID:     id,
			ActionBy:     actionBy,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Deleted member with ID " + id.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		span.AddEvent(`members.svc.delete.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionDeleted,
			ActionType:   "delete_member",
			ActionID:     id,
			ActionBy:     actionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Delete member failed: %v", err),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`members.svc.delete.success`)
	return nil
}
