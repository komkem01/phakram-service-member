package members

import (
	"context"
	"fmt"
	"strings"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/hashing"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdatePasswordServiceRequest struct {
	Password string
	ActionBy *uuid.UUID
}

func (s *Service) UpdatePasswordService(ctx context.Context, id uuid.UUID, req *UpdatePasswordServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.update_password.start`)

	now := time.Now()
	passwordHash, err := hashing.HashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_member_password",
			ActionID:     id,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update member password failed: %v", err),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member := new(ent.MemberEntity)
		if err := tx.NewSelect().Model(member).Where("id = ?", id).Where("deleted_at IS NULL").Scan(ctx); err != nil {
			return err
		}

		res, err := tx.NewUpdate().
			Table("member_accounts").
			Set("password = ?", string(passwordHash)).
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
			Details:   "member password updated",
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_member_password",
			ActionID:     id,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated member password with ID " + id.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		span.AddEvent(`members.svc.update_password.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_member_password",
			ActionID:     id,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update member password failed: %v", err),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`members.svc.update_password.success`)
	return nil
}
