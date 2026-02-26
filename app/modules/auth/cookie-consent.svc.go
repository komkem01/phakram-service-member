package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const cookiePolicyKey = "cookie_notice"

type CookiePolicyInfoResponse struct {
	PolicyID       uuid.UUID `json:"policy_id"`
	PolicyKey      string    `json:"policy_key"`
	VersionNo      int       `json:"version_no"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	EffectiveAt    int64     `json:"effective_at"`
	Accepted       bool      `json:"accepted"`
	AcceptedAt     *int64    `json:"accepted_at"`
	RequireConsent bool      `json:"require_consent"`
}

type AcceptCookieConsentServiceRequest struct {
	VisitorKey string
	UserAgent  string
}

type CreateCookiePolicyVersionServiceRequest struct {
	Title   string
	Content string
}

type CookiePolicyVersionItem struct {
	PolicyID    uuid.UUID  `json:"policy_id"`
	PolicyKey   string     `json:"policy_key"`
	VersionNo   int        `json:"version_no"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	IsActive    bool       `json:"is_active"`
	EffectiveAt int64      `json:"effective_at"`
	CreatedAt   int64      `json:"created_at"`
	UpdatedAt   int64      `json:"updated_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
}

type cookiePolicyVersionRow struct {
	ID          uuid.UUID  `bun:"id"`
	PolicyKey   string     `bun:"policy_key"`
	VersionNo   int        `bun:"version_no"`
	Title       string     `bun:"title"`
	Content     string     `bun:"content"`
	EffectiveAt *time.Time `bun:"effective_at"`
}

type cookiePolicyConsentRow struct {
	AcceptedAt *time.Time `bun:"accepted_at"`
}

type cookiePolicyConsentEntity struct {
	ID              uuid.UUID  `bun:"id,pk,type:uuid"`
	PolicyVersionID uuid.UUID  `bun:"policy_version_id,type:uuid"`
	MemberID        *uuid.UUID `bun:"member_id,type:uuid"`
	VisitorKey      string     `bun:"visitor_key"`
	AcceptedAt      time.Time  `bun:"accepted_at"`
	UserAgent       string     `bun:"user_agent"`
	CreatedAt       time.Time  `bun:"created_at"`
	UpdatedAt       time.Time  `bun:"updated_at"`
}

type cookiePolicyVersionEntity struct {
	ID          uuid.UUID  `bun:"id,pk,type:uuid"`
	PolicyKey   string     `bun:"policy_key"`
	VersionNo   int        `bun:"version_no"`
	Title       string     `bun:"title"`
	Content     string     `bun:"content"`
	IsActive    bool       `bun:"is_active"`
	EffectiveAt time.Time  `bun:"effective_at"`
	CreatedBy   *uuid.UUID `bun:"created_by,type:uuid"`
	CreatedAt   time.Time  `bun:"created_at"`
	UpdatedAt   time.Time  `bun:"updated_at"`
}

func mapCookiePolicyVersionEntity(item *cookiePolicyVersionEntity) *CookiePolicyVersionItem {
	if item == nil {
		return nil
	}
	return &CookiePolicyVersionItem{
		PolicyID:    item.ID,
		PolicyKey:   item.PolicyKey,
		VersionNo:   item.VersionNo,
		Title:       item.Title,
		Content:     item.Content,
		IsActive:    item.IsActive,
		EffectiveAt: item.EffectiveAt.Unix(),
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
		CreatedBy:   item.CreatedBy,
	}
}

func normalizeCookieVisitorKey(input string, memberID *uuid.UUID) string {
	if memberID != nil && *memberID != uuid.Nil {
		return "member:" + memberID.String()
	}
	return strings.TrimSpace(input)
}

func (s *Service) getActiveCookiePolicyVersion(ctx context.Context) (*cookiePolicyVersionRow, error) {
	row := new(cookiePolicyVersionRow)
	err := s.bunDB.DB().NewSelect().
		Table("cookie_policy_versions").
		Column("id", "policy_key", "version_no", "title", "content", "effective_at").
		Where("policy_key = ?", cookiePolicyKey).
		Where("is_active = ?", true).
		OrderExpr("version_no DESC").
		OrderExpr("effective_at DESC").
		Limit(1).
		Scan(ctx, row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cookie policy not configured")
		}
		return nil, err
	}
	return row, nil
}

func (s *Service) getCookieConsentAcceptedAt(ctx context.Context, policyVersionID uuid.UUID, visitorKey string) (*time.Time, error) {
	if strings.TrimSpace(visitorKey) == "" {
		return nil, nil
	}

	row := new(cookiePolicyConsentRow)
	err := s.bunDB.DB().NewSelect().
		Table("cookie_policy_consents").
		Column("accepted_at").
		Where("policy_version_id = ?", policyVersionID).
		Where("visitor_key = ?", visitorKey).
		Limit(1).
		Scan(ctx, row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return row.AcceptedAt, nil
}

func (s *Service) GetCookiePolicyInfoService(ctx context.Context, memberID *uuid.UUID, visitorKey string) (*CookiePolicyInfoResponse, error) {
	policy, err := s.getActiveCookiePolicyVersion(ctx)
	if err != nil {
		return nil, err
	}

	resolvedVisitorKey := normalizeCookieVisitorKey(visitorKey, memberID)
	acceptedAt, err := s.getCookieConsentAcceptedAt(ctx, policy.ID, resolvedVisitorKey)
	if err != nil {
		return nil, err
	}

	response := &CookiePolicyInfoResponse{
		PolicyID:       policy.ID,
		PolicyKey:      policy.PolicyKey,
		VersionNo:      policy.VersionNo,
		Title:          policy.Title,
		Content:        policy.Content,
		EffectiveAt:    time.Now().Unix(),
		Accepted:       acceptedAt != nil,
		RequireConsent: acceptedAt == nil,
	}
	if policy.EffectiveAt != nil {
		response.EffectiveAt = policy.EffectiveAt.Unix()
	}
	if acceptedAt != nil {
		acceptedUnix := acceptedAt.Unix()
		response.AcceptedAt = &acceptedUnix
	}

	return response, nil
}

func (s *Service) AcceptCookieConsentService(ctx context.Context, req *AcceptCookieConsentServiceRequest, memberID *uuid.UUID) (*CookiePolicyInfoResponse, error) {
	policy, err := s.getActiveCookiePolicyVersion(ctx)
	if err != nil {
		return nil, err
	}

	visitorKey := normalizeCookieVisitorKey(req.VisitorKey, memberID)
	if strings.TrimSpace(visitorKey) == "" {
		return nil, errors.New("visitor key is required")
	}

	now := time.Now()
	consent := &cookiePolicyConsentEntity{
		ID:              uuid.New(),
		PolicyVersionID: policy.ID,
		MemberID:        memberID,
		VisitorKey:      visitorKey,
		AcceptedAt:      now,
		UserAgent:       strings.TrimSpace(req.UserAgent),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err = s.bunDB.DB().NewInsert().
		Table("cookie_policy_consents").
		Model(consent).
		On("CONFLICT (policy_version_id, visitor_key) DO UPDATE").
		Set("member_id = EXCLUDED.member_id").
		Set("accepted_at = EXCLUDED.accepted_at").
		Set("user_agent = EXCLUDED.user_agent").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return s.GetCookiePolicyInfoService(ctx, memberID, visitorKey)
}

func (s *Service) ListCookiePolicyVersionsService(ctx context.Context) ([]*CookiePolicyVersionItem, error) {
	rows := make([]*cookiePolicyVersionEntity, 0)
	err := s.bunDB.DB().NewSelect().
		Table("cookie_policy_versions").
		Column("id", "policy_key", "version_no", "title", "content", "is_active", "effective_at", "created_by", "created_at", "updated_at").
		Where("policy_key = ?", cookiePolicyKey).
		OrderExpr("version_no DESC").
		Scan(ctx, &rows)
	if err != nil {
		return nil, err
	}

	items := make([]*CookiePolicyVersionItem, 0, len(rows))
	for _, row := range rows {
		mapped := mapCookiePolicyVersionEntity(row)
		if mapped != nil {
			items = append(items, mapped)
		}
	}

	return items, nil
}

func (s *Service) CreateCookiePolicyVersionService(ctx context.Context, req *CreateCookiePolicyVersionServiceRequest, createdBy uuid.UUID) (*CookiePolicyVersionItem, error) {
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)

	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}

	now := time.Now()
	createdByID := createdBy
	newVersion := &cookiePolicyVersionEntity{}

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		latest := new(cookiePolicyVersionEntity)
		err := tx.NewSelect().
			Table("cookie_policy_versions").
			Column("id", "policy_key", "version_no", "title", "content", "is_active", "effective_at", "created_by", "created_at", "updated_at").
			Where("policy_key = ?", cookiePolicyKey).
			OrderExpr("version_no DESC").
			Limit(1).
			Scan(ctx, latest)

		nextVersion := 1
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		} else {
			nextVersion = latest.VersionNo + 1
		}

		if _, err := tx.NewUpdate().
			Table("cookie_policy_versions").
			Set("is_active = ?", false).
			Set("updated_at = ?", now).
			Where("policy_key = ?", cookiePolicyKey).
			Where("is_active = ?", true).
			Exec(ctx); err != nil {
			return err
		}

		newVersion = &cookiePolicyVersionEntity{
			ID:          uuid.New(),
			PolicyKey:   cookiePolicyKey,
			VersionNo:   nextVersion,
			Title:       title,
			Content:     content,
			IsActive:    true,
			EffectiveAt: now,
			CreatedBy:   &createdByID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if _, err := tx.NewInsert().
			Table("cookie_policy_versions").
			Model(newVersion).
			Exec(ctx); err != nil {
			return fmt.Errorf("create cookie policy version failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapCookiePolicyVersionEntity(newVersion), nil
}
