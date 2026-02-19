package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type authSessionEntity struct {
	bun.BaseModel    `bun:"table:auth_sessions"`
	ID               uuid.UUID  `bun:"id,pk,type:uuid"`
	MemberID         uuid.UUID  `bun:"member_id,type:uuid"`
	ActorMemberID    *uuid.UUID `bun:"actor_member_id,type:uuid"`
	ActorIsAdmin     bool       `bun:"actor_is_admin"`
	IsActingAs       bool       `bun:"is_acting_as"`
	RefreshTokenHash string     `bun:"refresh_token_hash"`
	LastActivity     time.Time  `bun:"last_activity"`
	RefreshExpiresAt time.Time  `bun:"refresh_expires_at"`
	RevokedAt        *time.Time `bun:"revoked_at"`
	CreatedAt        time.Time  `bun:"created_at"`
	UpdatedAt        time.Time  `bun:"updated_at"`
}

func hashTokenValue(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *Service) createAuthSession(ctx context.Context, sessionID uuid.UUID, memberID uuid.UUID, actorMemberID *uuid.UUID, actorIsAdmin bool, isActingAs bool, refreshToken string, refreshExpUnix int64) (*authSessionEntity, error) {
	now := time.Now().UTC()
	item := &authSessionEntity{
		ID:               sessionID,
		MemberID:         memberID,
		ActorMemberID:    actorMemberID,
		ActorIsAdmin:     actorIsAdmin,
		IsActingAs:       isActingAs,
		RefreshTokenHash: hashTokenValue(refreshToken),
		LastActivity:     now,
		RefreshExpiresAt: time.Unix(refreshExpUnix, 0).UTC(),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if _, err := s.bunDB.DB().NewInsert().Model(item).Exec(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) getAuthSessionByID(ctx context.Context, sessionID uuid.UUID) (*authSessionEntity, error) {
	item := new(authSessionEntity)
	err := s.bunDB.DB().NewSelect().
		Model(item).
		Where("id = ?", sessionID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) isSessionExpiredByIdle(lastActivity time.Time, now time.Time) bool {
	return now.Sub(lastActivity) > idleSessionTTL
}

func (s *Service) validateAccessSession(ctx context.Context, sessionID uuid.UUID, memberID uuid.UUID) (*authSessionEntity, error) {
	item, err := s.getAuthSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if item.MemberID != memberID {
		return nil, sql.ErrNoRows
	}
	if item.RevokedAt != nil {
		return nil, sql.ErrNoRows
	}
	if now.After(item.RefreshExpiresAt) {
		return nil, sql.ErrNoRows
	}
	if s.isSessionExpiredByIdle(item.LastActivity, now) {
		return nil, sql.ErrNoRows
	}

	return item, nil
}

func (s *Service) touchAuthSessionActivity(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now().UTC()
	_, err := s.bunDB.DB().NewUpdate().
		Table("auth_sessions").
		Set("last_activity = ?", now).
		Set("updated_at = ?", now).
		Where("id = ?", sessionID).
		Where("revoked_at IS NULL").
		Exec(ctx)
	return err
}

func (s *Service) rotateRefreshSession(ctx context.Context, sessionID uuid.UUID, refreshToken string, refreshExpUnix int64) error {
	now := time.Now().UTC()
	_, err := s.bunDB.DB().NewUpdate().
		Table("auth_sessions").
		Set("refresh_token_hash = ?", hashTokenValue(refreshToken)).
		Set("refresh_expires_at = ?", time.Unix(refreshExpUnix, 0).UTC()).
		Set("last_activity = ?", now).
		Set("updated_at = ?", now).
		Where("id = ?", sessionID).
		Where("revoked_at IS NULL").
		Exec(ctx)
	return err
}

func (s *Service) revokeSession(ctx context.Context, sessionID uuid.UUID, memberID uuid.UUID) error {
	now := time.Now().UTC()
	_, err := s.bunDB.DB().NewUpdate().
		Table("auth_sessions").
		Set("revoked_at = ?", now).
		Set("updated_at = ?", now).
		Where("id = ?", sessionID).
		Where("member_id = ?", memberID).
		Where("revoked_at IS NULL").
		Exec(ctx)
	return err
}
