package member_accounts

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/hashing"

	"github.com/google/uuid"
)

type CreateMemberAccountService struct {
	MemberID uuid.UUID `json:"member_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (s *Service) CreateMemberAccountService(ctx context.Context, req *CreateMemberAccountService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_accounts.svc.create.start`)

	hashedPassword, err := hashing.HashPassword(req.Password)
	if err != nil {
		return err
	}

	account := &ent.MemberAccountEntity{
		ID:       uuid.New(),
		MemberID: req.MemberID,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	if err := s.db.CreateMemberAccount(ctx, account); err != nil {
		return err
	}
	span.AddEvent(`member_accounts.svc.create.success`)
	return nil
}
