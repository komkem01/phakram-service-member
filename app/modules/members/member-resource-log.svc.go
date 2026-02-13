package members

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"phakram/app/modules/auth"
	"phakram/app/modules/entities/ent"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type auditActionDetail struct {
	Message       string `json:"message,omitempty"`
	Endpoint      string `json:"endpoint,omitempty"`
	ActorMemberID string `json:"actor_member_id,omitempty"`
	TargetMemberID string `json:"target_member_id,omitempty"`
	ActingAs      bool   `json:"acting_as,omitempty"`
	Before        any    `json:"before,omitempty"`
	After         any    `json:"after,omitempty"`
	Error         string `json:"error,omitempty"`
}

func (s *Service) buildAuditActionDetail(
	ctx context.Context,
	memberID uuid.UUID,
	actionBy *uuid.UUID,
	message string,
	before any,
	after any,
	err error,
) string {
	actorMemberID, hasActor := auth.RequestActorMemberID(ctx)
	if !hasActor && actionBy != nil {
		actorMemberID = *actionBy
		hasActor = true
	}

	targetMemberID, hasTarget := auth.RequestTargetMemberID(ctx)
	if !hasTarget && memberID != uuid.Nil {
		targetMemberID = memberID
		hasTarget = true
	}

	detail := auditActionDetail{
		Message:  message,
		Endpoint: auth.RequestEndpoint(ctx),
		ActingAs: auth.RequestIsActingAs(ctx),
		Before:   before,
		After:    after,
	}
	if hasActor {
		detail.ActorMemberID = actorMemberID.String()
	}
	if hasTarget {
		detail.TargetMemberID = targetMemberID.String()
	}
	if err != nil {
		detail.Error = err.Error()
	}

	payload, marshalErr := json.Marshal(detail)
	if marshalErr != nil {
		if err != nil {
			return fmt.Sprintf("%s failed: %v", message, err)
		}
		return message
	}

	return string(payload)
}

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
		ActionDetail: s.buildAuditActionDetail(ctx, memberID, actionBy, detail, nil, nil, nil),
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
		ActionDetail: s.buildAuditActionDetail(ctx, actionID, actionBy, actionType, nil, nil, err),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
}
