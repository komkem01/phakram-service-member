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

type UpdateServiceRequest struct {
	TierID      *uuid.UUID
	StatusID    *uuid.UUID
	PrefixID    *uuid.UUID
	GenderID    *uuid.UUID
	FirstnameTh *string
	LastnameTh  *string
	FirstnameEn *string
	LastnameEn  *string
	Role        *string
	Phone       *string
	ActionBy    *uuid.UUID
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.update.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", id).Where("deleted_at IS NULL").Scan(ctx); err != nil {
			return err
		}

		if req.TierID != nil {
			data.TierID = *req.TierID
		}
		if req.StatusID != nil {
			data.StatusID = *req.StatusID
		}
		if req.PrefixID != nil {
			data.PrefixID = *req.PrefixID
		}
		if req.GenderID != nil {
			data.GenderID = *req.GenderID
		}
		if req.FirstnameTh != nil {
			data.FirstnameTh = strings.TrimSpace(*req.FirstnameTh)
		}
		if req.LastnameTh != nil {
			data.LastnameTh = strings.TrimSpace(*req.LastnameTh)
		}
		if req.FirstnameEn != nil {
			data.FirstnameEn = strings.TrimSpace(*req.FirstnameEn)
		}
		if req.LastnameEn != nil {
			data.LastnameEn = strings.TrimSpace(*req.LastnameEn)
		}
		if req.Phone != nil {
			data.Phone = strings.TrimSpace(*req.Phone)
		}
		if req.Role != nil {
			switch strings.TrimSpace(*req.Role) {
			case string(ent.RoleTypeAdmin):
				data.Role = ent.RoleTypeAdmin
			case string(ent.RoleTypeCustomer):
				data.Role = ent.RoleTypeCustomer
			}
		}
		data.UpdatedAt = now

		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		memberTx := &ent.MemberTransactionEntity{
			ID:        uuid.New(),
			MemberID:  id,
			Action:    ent.MemberActionUpdated,
			Details:   "member updated",
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_member",
			ActionID:     id,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated member with ID " + id.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		span.AddEvent(`members.svc.update.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_member",
			ActionID:     id,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update member failed: %v", err),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`members.svc.update.success`)
	return nil
}
