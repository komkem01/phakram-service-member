package members

import (
	"context"
	"fmt"
	"strings"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdateEmailServiceRequest struct {
	Email    string
	ActionBy *uuid.UUID
}

func (s *Service) UpdateEmailService(ctx context.Context, id uuid.UUID, req *UpdateEmailServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.update_email.start`)

	now := time.Now()
	email := strings.TrimSpace(req.Email)

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member := new(ent.MemberEntity)
		if err := tx.NewSelect().Model(member).Where("id = ?", id).Where("deleted_at IS NULL").Scan(ctx); err != nil {
			return err
		}

		res, err := tx.NewUpdate().
			Table("member_accounts").
			Set("email = ?", email).
			Set("updated_at = ?", now).
			Where("member_id = ?", id).
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
			Action:    ent.MemberActionUpdated,
			Details:   "member email updated",
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_member_email",
			ActionID:     id,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated member email with ID " + id.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		span.AddEvent(`members.svc.update_email.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_member_email",
			ActionID:     id,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update member email failed: %v", err),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`members.svc.update_email.success`)
	return nil
}
