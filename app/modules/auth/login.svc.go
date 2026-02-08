package auth

import (
	"context"
	"database/sql"
	"errors"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/hashing"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginMemberService struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginMemberServiceResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) LoginMemberService(ctx context.Context, req *LoginMemberService) (*LoginMemberServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.login.start`)

	account, err := s.dbAccount.GetMemberAccountByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !hashing.CheckPasswordHash([]byte(account.Password), []byte(req.Password)) {
		return nil, ErrInvalidCredentials
	}

	member, err := s.db.GetMemberByID(ctx, account.MemberID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	roleLabel := "member"
	if member.Role == ent.RoleTypeAdmin {
		roleLabel = "admin"
	}

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member.LastLogin = &now
		if _, err := tx.NewUpdate().
			Model(member).
			Column("last_login").
			Where("id = ?", member.ID).
			Exec(ctx); err != nil {
			return err
		}

		memberTransaction := &ent.MemberTransactionEntity{
			ID:        uuid.New(),
			MemberID:  member.ID,
			Action:    ent.ActionTypeLogined,
			Details:   "login: " + roleLabel,
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(memberTransaction).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	isAdmin := member.Role == ent.RoleTypeAdmin
	accessToken, _, err := s.generateToken(member.ID, account.Email, string(member.Role), isAdmin, "access", accessTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, _, err := s.generateToken(member.ID, account.Email, string(member.Role), isAdmin, "refresh", refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	resp := &LoginMemberServiceResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	span.AddEvent(`members.svc.login.success`)
	return resp, nil
}
