package members

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/hashing"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RegisterServiceRequest struct {
	PrefixID    uuid.UUID
	GenderID    uuid.UUID
	FirstnameTh string
	LastnameTh  string
	FirstnameEn string
	LastnameEn  string
	Role        string
	Phone       string
	Email       string
	Password    string
	TierID      uuid.UUID
	StatusID    uuid.UUID
	ActionBy    *uuid.UUID
}

func (s *Service) CreateRegisterService(ctx context.Context, req *RegisterServiceRequest) error {
	return s.createRegisterService(ctx, req, false)
}

func (s *Service) CreateRegisterByAdminService(ctx context.Context, req *RegisterServiceRequest) error {
	return s.createRegisterService(ctx, req, true)
}

func (s *Service) createRegisterService(ctx context.Context, req *RegisterServiceRequest, allowRoleSelection bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`register.svc.create.start`)

	memberID := uuid.New()
	accountID := uuid.New()

	passwordHash, err := hashing.HashPassword(req.Password)
	if err != nil {
		span.AddEvent(`register.svc.create.hash_failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionRegistered,
			ActionType:   s.registerActionType(allowRoleSelection),
			ActionID:     memberID,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: s.buildAuditActionDetail(ctx, memberID, req.ActionBy, "Register member", nil, nil, err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	role := ent.RoleTypeCustomer
	if allowRoleSelection {
		switch strings.ToLower(strings.TrimSpace(req.Role)) {
		case string(ent.RoleTypeAdmin):
			role = ent.RoleTypeAdmin
		case string(ent.RoleTypeCustomer):
			role = ent.RoleTypeCustomer
		default:
			return errors.New("invalid role")
		}
	}

	now := time.Now()

	tierID, statusID, err := s.resolveTierStatus(ctx, req, allowRoleSelection)
	if err != nil {
		span.AddEvent(`register.svc.create.resolve_tier_status_failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionRegistered,
			ActionType:   s.registerActionType(allowRoleSelection),
			ActionID:     memberID,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: s.buildAuditActionDetail(ctx, memberID, req.ActionBy, "Register member", nil, nil, err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	member := &ent.MemberEntity{
		ID:           memberID,
		PrefixID:     req.PrefixID,
		GenderID:     req.GenderID,
		TierID:       tierID,
		StatusID:     statusID,
		FirstnameTh:  req.FirstnameTh,
		LastnameTh:   req.LastnameTh,
		FirstnameEn:  req.FirstnameEn,
		LastnameEn:   req.LastnameEn,
		Role:         role,
		Phone:        req.Phone,
		Registration: &now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	memberAccount := &ent.MemberAccountEntity{
		ID:        accountID,
		MemberID:  memberID,
		Email:     req.Email,
		Password:  string(passwordHash),
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member.MemberNo = utils.GenerateMemberNo(memberID, string(role))

		if _, err := tx.NewInsert().Model(member).Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewInsert().Model(memberAccount).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionRegistered,
			ActionType:   s.registerActionType(allowRoleSelection),
			ActionID:     memberID,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: s.buildAuditActionDetail(ctx, memberID, req.ActionBy, "Registered member", nil, nil, nil),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		memberTx := &ent.MemberTransactionEntity{
			ID:        uuid.New(),
			MemberID:  memberID,
			Action:    ent.MemberActionRegistered,
			Details:   "member registered",
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		span.AddEvent(`register.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionRegistered,
			ActionType:   s.registerActionType(allowRoleSelection),
			ActionID:     memberID,
			ActionBy:     req.ActionBy,
			Status:       ent.StatusAuditFailed,
			ActionDetail: s.buildAuditActionDetail(ctx, memberID, req.ActionBy, "Register member", nil, nil, err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`register.svc.create.success`)
	return nil
}

func (s *Service) registerActionType(allowRoleSelection bool) string {
	if allowRoleSelection {
		return "register_member_admin"
	}
	return "register_member"
}

func (s *Service) resolveTierStatus(ctx context.Context, req *RegisterServiceRequest, allowRoleSelection bool) (uuid.UUID, uuid.UUID, error) {
	if allowRoleSelection && (req.TierID != uuid.Nil || req.StatusID != uuid.Nil) {
		tierID := req.TierID
		statusID := req.StatusID

		if tierID != uuid.Nil {
			if err := s.ensureTierExists(ctx, tierID); err != nil {
				return uuid.Nil, uuid.Nil, err
			}
		} else {
			defaultTierID, err := s.defaultTierID(ctx)
			if err != nil {
				return uuid.Nil, uuid.Nil, err
			}
			tierID = defaultTierID
		}

		if statusID != uuid.Nil {
			if err := s.ensureStatusExists(ctx, statusID); err != nil {
				return uuid.Nil, uuid.Nil, err
			}
		} else {
			defaultStatusID, err := s.defaultStatusID(ctx)
			if err != nil {
				return uuid.Nil, uuid.Nil, err
			}
			statusID = defaultStatusID
		}

		return tierID, statusID, nil
	}

	defaultTierID, err := s.defaultTierID(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	defaultStatusID, err := s.defaultStatusID(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return defaultTierID, defaultStatusID, nil
}

func (s *Service) defaultTierID(ctx context.Context) (uuid.UUID, error) {
	data := new(ent.TierEntity)
	err := s.bunDB.DB().NewSelect().
		Model(data).
		Where("is_active = true").
		OrderExpr("min_spending ASC, created_at ASC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, errors.New("default tier not found")
		}
		return uuid.Nil, err
	}
	return data.ID, nil
}

func (s *Service) defaultStatusID(ctx context.Context) (uuid.UUID, error) {
	data := new(ent.StatusEntity)
	err := s.bunDB.DB().NewSelect().
		Model(data).
		Where("is_active = true").
		OrderExpr("created_at ASC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, errors.New("default status not found")
		}
		return uuid.Nil, err
	}
	return data.ID, nil
}

func (s *Service) ensureTierExists(ctx context.Context, id uuid.UUID) error {
	data := new(ent.TierEntity)
	return s.bunDB.DB().NewSelect().
		Model(data).
		Column("id").
		Where("id = ?", id).
		Scan(ctx)
}

func (s *Service) ensureStatusExists(ctx context.Context, id uuid.UUID) error {
	data := new(ent.StatusEntity)
	return s.bunDB.DB().NewSelect().
		Model(data).
		Column("id").
		Where("id = ?", id).
		Scan(ctx)
}
