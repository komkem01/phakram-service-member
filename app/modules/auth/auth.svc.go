package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/hashing"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type LoginServiceRequest struct {
	Email    string
	Password string
}

type TokenServiceResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	AccessExpiresAt  int64  `json:"access_expires_at"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
}

type RefreshTokenServiceRequest struct {
	RefreshToken string
}

type RevokeSessionServiceRequest struct {
	MemberID  uuid.UUID
	SessionID uuid.UUID
}

type ActAsMemberServiceRequest struct {
	ActorMemberID  uuid.UUID
	TargetMemberID uuid.UUID
}

type MemberInfoServiceResponse struct {
	MemberID      uuid.UUID  `json:"member_id"`
	MemberNo      string     `json:"member_no"`
	Email         string     `json:"email"`
	Role          string     `json:"role"`
	IsAdmin       bool       `json:"is_admin"`
	FirstnameTh   string     `json:"firstname_th"`
	LastnameTh    string     `json:"lastname_th"`
	FirstnameEn   string     `json:"firstname_en"`
	LastnameEn    string     `json:"lastname_en"`
	Phone         string     `json:"phone"`
	ActorMemberID *uuid.UUID `json:"actor_member_id,omitempty"`
	ActorIsAdmin  bool       `json:"actor_is_admin"`
	IsActingAs    bool       `json:"is_acting_as"`
	LastLogin     *int64     `json:"last_login"`
	Registration  *int64     `json:"registration"`
}

type memberAuthData struct {
	MemberID     uuid.UUID        `bun:"member_id"`
	MemberNo     string           `bun:"member_no"`
	Email        string           `bun:"email"`
	Password     string           `bun:"password"`
	Role         ent.RoleTypeEnum `bun:"role"`
	FirstnameTh  string           `bun:"firstname_th"`
	LastnameTh   string           `bun:"lastname_th"`
	FirstnameEn  string           `bun:"firstname_en"`
	LastnameEn   string           `bun:"lastname_en"`
	Phone        string           `bun:"phone"`
	LastLogin    *time.Time       `bun:"last_login"`
	Registration *time.Time       `bun:"registration"`
}

func (s *Service) LoginService(ctx context.Context, req *LoginServiceRequest) (*TokenServiceResponse, error) {
	_, span, _ := utils.NewLogSpan(ctx, s.tracer, "LoginService")
	defer span.End()

	data, err := s.findMemberAuthByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !hashing.CheckPasswordHash([]byte(data.Password), []byte(req.Password)) {
		return nil, errors.New("invalid credentials")
	}

	if err := s.touchMemberLastLogin(ctx, data.MemberID, fmt.Sprintf("member logined by email %s", data.Email)); err != nil {
		return nil, err
	}

	return s.buildTokenResponse(ctx, data.MemberID, data.Email, data.Role)
}

func (s *Service) RefreshTokenService(ctx context.Context, req *RefreshTokenServiceRequest) (*TokenServiceResponse, error) {
	_, span, _ := utils.NewLogSpan(ctx, s.tracer, "RefreshTokenService")
	defer span.End()

	claims, err := s.parseToken(req.RefreshToken, "refresh")
	if err != nil {
		return nil, errors.New("invalid token")
	}

	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	memberID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	session, err := s.getAuthSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	now := time.Now().UTC()
	if session.MemberID != memberID || session.RevokedAt != nil || now.After(session.RefreshExpiresAt) || s.isSessionExpiredByIdle(session.LastActivity, now) {
		return nil, errors.New("invalid token")
	}
	if session.RefreshTokenHash != hashTokenValue(req.RefreshToken) {
		return nil, errors.New("invalid token")
	}

	data, err := s.findMemberAuthByID(ctx, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	if err := s.touchMemberLastLogin(ctx, data.MemberID, "member refreshed token"); err != nil {
		return nil, err
	}

	actorSub := ""
	if session.ActorMemberID != nil {
		actorSub = session.ActorMemberID.String()
	}

	res, err := s.buildTokenResponseWithSession(data.MemberID, data.Email, data.Role, session.ID, actorSub, session.ActorIsAdmin, session.IsActingAs)
	if err != nil {
		return nil, err
	}

	if err := s.rotateRefreshSession(ctx, session.ID, res.RefreshToken, res.RefreshExpiresAt); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) RevokeSessionService(ctx context.Context, req *RevokeSessionServiceRequest) error {
	_, span, _ := utils.NewLogSpan(ctx, s.tracer, "RevokeSessionService")
	defer span.End()

	return s.revokeSession(ctx, req.SessionID, req.MemberID)
}

func (s *Service) touchMemberLastLogin(ctx context.Context, memberID uuid.UUID, detail string) error {
	now := time.Now()

	return s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().
			Table("members").
			Set("last_login = ?", now).
			Set("updated_at = ?", now).
			Where("id = ?", memberID).
			Exec(ctx); err != nil {
			return err
		}

		memberTx := &ent.MemberTransactionEntity{
			ID:        uuid.New(),
			MemberID:  memberID,
			Action:    ent.MemberActionLogined,
			Details:   detail,
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) GetInfoService(ctx context.Context, memberID uuid.UUID) (*MemberInfoServiceResponse, error) {
	_, span, _ := utils.NewLogSpan(ctx, s.tracer, "GetInfoService")
	defer span.End()

	data, err := s.findMemberAuthByID(ctx, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	resp := &MemberInfoServiceResponse{
		MemberID:     data.MemberID,
		MemberNo:     data.MemberNo,
		Email:        data.Email,
		Role:         string(data.Role),
		IsAdmin:      data.Role == ent.RoleTypeAdmin,
		FirstnameTh:  data.FirstnameTh,
		LastnameTh:   data.LastnameTh,
		FirstnameEn:  data.FirstnameEn,
		LastnameEn:   data.LastnameEn,
		Phone:        data.Phone,
		ActorIsAdmin: RequestIsActingAs(ctx) || (data.Role == ent.RoleTypeAdmin),
		IsActingAs:   RequestIsActingAs(ctx),
	}

	if actorID, ok := RequestActorMemberID(ctx); ok && actorID != memberID {
		resp.ActorMemberID = &actorID
		resp.ActorIsAdmin = true
	}

	if data.LastLogin != nil {
		lastLogin := data.LastLogin.Unix()
		resp.LastLogin = &lastLogin
	}
	if data.Registration != nil {
		registration := data.Registration.Unix()
		resp.Registration = &registration
	}

	return resp, nil
}

func (s *Service) ActAsMemberService(ctx context.Context, req *ActAsMemberServiceRequest) (*TokenServiceResponse, error) {
	_, span, _ := utils.NewLogSpan(ctx, s.tracer, "ActAsMemberService")
	defer span.End()

	actorData, err := s.findMemberAuthByID(ctx, req.ActorMemberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("forbidden")
		}
		return nil, err
	}
	if actorData.Role != ent.RoleTypeAdmin {
		return nil, errors.New("forbidden")
	}

	targetData, err := s.findMemberAuthByID(ctx, req.TargetMemberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("forbidden")
		}
		return nil, err
	}

	return s.buildTokenResponseWithActor(ctx, targetData.MemberID, targetData.Email, targetData.Role, actorData.MemberID.String(), true, true)
}

func (s *Service) ExitActAsService(ctx context.Context, actorMemberID uuid.UUID) (*TokenServiceResponse, error) {
	_, span, _ := utils.NewLogSpan(ctx, s.tracer, "ExitActAsService")
	defer span.End()

	actorData, err := s.findMemberAuthByID(ctx, actorMemberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("forbidden")
		}
		return nil, err
	}
	if actorData.Role != ent.RoleTypeAdmin {
		return nil, errors.New("forbidden")
	}

	return s.buildTokenResponse(ctx, actorData.MemberID, actorData.Email, actorData.Role)
}

func (s *Service) buildTokenResponse(ctx context.Context, memberID uuid.UUID, email string, role ent.RoleTypeEnum) (*TokenServiceResponse, error) {
	return s.buildTokenResponseWithActor(ctx, memberID, email, role, "", false, false)
}

func (s *Service) buildTokenResponseWithActor(ctx context.Context, memberID uuid.UUID, email string, role ent.RoleTypeEnum, actorSub string, actorIsAdmin bool, actingAs bool) (*TokenServiceResponse, error) {
	sessionID := uuid.New()

	res, err := s.buildTokenResponseWithSession(memberID, email, role, sessionID, actorSub, actorIsAdmin, actingAs)
	if err != nil {
		return nil, err
	}

	var actorMemberID *uuid.UUID
	if actorSub != "" {
		parsedActorMemberID, parseErr := uuid.Parse(actorSub)
		if parseErr != nil {
			return nil, parseErr
		}
		actorMemberID = &parsedActorMemberID
	}

	if _, err := s.createAuthSession(ctx, sessionID, memberID, actorMemberID, actorIsAdmin, actingAs, res.RefreshToken, res.RefreshExpiresAt); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) buildTokenResponseWithSession(memberID uuid.UUID, email string, role ent.RoleTypeEnum, sessionID uuid.UUID, actorSub string, actorIsAdmin bool, actingAs bool) (*TokenServiceResponse, error) {
	isAdmin := role == ent.RoleTypeAdmin

	accessToken, accessExp, err := s.generateToken(memberID, email, string(role), isAdmin, sessionID, actorSub, actorIsAdmin, actingAs, "access", accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExp, err := s.generateToken(memberID, email, string(role), isAdmin, sessionID, actorSub, actorIsAdmin, actingAs, "refresh", refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &TokenServiceResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}, nil
}

func (s *Service) findMemberAuthByEmail(ctx context.Context, email string) (*memberAuthData, error) {
	data := new(memberAuthData)
	err := s.bunDB.DB().NewSelect().
		TableExpr("member_accounts AS ma").
		ColumnExpr("ma.member_id").
		ColumnExpr("m.member_no").
		ColumnExpr("ma.email").
		ColumnExpr("ma.password").
		ColumnExpr("m.role").
		ColumnExpr("m.firstname_th").
		ColumnExpr("m.lastname_th").
		ColumnExpr("m.firstname_en").
		ColumnExpr("m.lastname_en").
		ColumnExpr("m.phone").
		ColumnExpr("m.last_login").
		ColumnExpr("m.registration").
		Join("JOIN members AS m ON m.id = ma.member_id").
		Where("ma.email = ?", email).
		Where("m.deleted_at IS NULL").
		Limit(1).
		Scan(ctx, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) findMemberAuthByID(ctx context.Context, memberID uuid.UUID) (*memberAuthData, error) {
	data := new(memberAuthData)
	err := s.bunDB.DB().NewSelect().
		TableExpr("members AS m").
		ColumnExpr("m.id AS member_id").
		ColumnExpr("m.member_no").
		ColumnExpr("ma.email").
		ColumnExpr("ma.password").
		ColumnExpr("m.role").
		ColumnExpr("m.firstname_th").
		ColumnExpr("m.lastname_th").
		ColumnExpr("m.firstname_en").
		ColumnExpr("m.lastname_en").
		ColumnExpr("m.phone").
		ColumnExpr("m.last_login").
		ColumnExpr("m.registration").
		Join("JOIN member_accounts AS ma ON ma.member_id = m.id").
		Where("m.id = ?", memberID).
		Where("m.deleted_at IS NULL").
		Limit(1).
		Scan(ctx, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
