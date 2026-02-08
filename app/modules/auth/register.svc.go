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
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type RegisterMemberService struct {
	MemberNo      string          `json:"member_no"`
	TierID        uuid.UUID       `json:"tier_id"`
	StatusID      uuid.UUID       `json:"status_id"`
	PrefixID      uuid.UUID       `json:"prefix_id"`
	GenderID      uuid.UUID       `json:"gender_id"`
	FirstnameTh   string          `json:"firstname_th"`
	LastnameTh    string          `json:"lastname_th"`
	FirstnameEn   string          `json:"firstname_en"`
	LastnameEn    string          `json:"lastname_en"`
	Role          string          `json:"role"`
	Phone         string          `json:"phone"`
	Email         string          `json:"email"`
	Password      string          `json:"password"`
	TotalSpent    decimal.Decimal `json:"total_spent"`
	CurrentPoints int             `json:"current_points"`
	Registration  time.Time       `json:"registration"`
	LastLogin     time.Time       `json:"last_login"`
}

var (
	defaultMemberTierID   = uuid.MustParse("b4ac1d4d-5779-4e78-aa77-6f174fe2f91d")
	defaultMemberStatusID = uuid.MustParse("5fdb7969-e19e-4219-8825-b71a1c5c8bf7")
)

func (s *Service) RegisterMemberService(ctx context.Context, req *RegisterMemberService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.register.start`)

	if req.Email != "" {
		_, err := s.dbAccount.GetMemberAccountByEmail(ctx, req.Email)
		if err == nil {
			return errors.New("email already exists")
		}
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	if req.Phone != "" {
		_, err := s.db.GetMemberByPhone(ctx, req.Phone)
		if err == nil {
			return errors.New("phone already exists")
		}
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	if req.Role != string(ent.RoleTypeAdmin) {
		req.TierID = defaultMemberTierID
		req.StatusID = defaultMemberStatusID
	}

	memberNo := req.MemberNo
	if memberNo == "" {
		var err error
		memberNo, err = utils.GenerateMemberNo(ctx, s.bunDB.DB())
		if err != nil {
			return err
		}
	}
	passwordHash, err := hashing.HashPassword(req.Password)
	if err != nil {
		return err
	}

	roleLabel := "member"
	if req.Role == string(ent.RoleTypeAdmin) {
		roleLabel = "admin"
	}

	id := uuid.New()

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member := &ent.MemberEntity{
			ID:            id,
			MemberNo:      memberNo,
			TierID:        req.TierID,
			StatusID:      req.StatusID,
			PrefixID:      req.PrefixID,
			GenderID:      req.GenderID,
			FirstnameTh:   req.FirstnameTh,
			LastnameTh:    req.LastnameTh,
			FirstnameEn:   req.FirstnameEn,
			LastnameEn:    req.LastnameEn,
			Role:          ent.RoleTypeEnum(req.Role),
			Phone:         req.Phone,
			TotalSpent:    req.TotalSpent,
			CurrentPoints: req.CurrentPoints,
			Registration: func() *time.Time {
				now := time.Now()
				return &now
			}(),
			LastLogin: nil,
		}
		if _, err := tx.NewInsert().Model(member).Exec(ctx); err != nil {
			return err
		}

		memberAccount := &ent.MemberAccountEntity{
			ID:       uuid.New(),
			MemberID: id,
			Email:    req.Email,
			Password: string(passwordHash),
		}
		if _, err := tx.NewInsert().Model(memberAccount).Exec(ctx); err != nil {
			return err
		}

		memberTransaction := &ent.MemberTransactionEntity{
			ID:        uuid.New(),
			MemberID:  id,
			Action:    ent.ActionTypeRegistered,
			Details:   "register: " + roleLabel,
			CreatedAt: time.Now(),
		}
		if _, err := tx.NewInsert().Model(memberTransaction).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	span.AddEvent(`members.svc.register.success`)
	return nil
}
