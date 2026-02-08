package member_accounts

import (
	"context"
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/hashing"

	"github.com/google/uuid"
)

type UpdateMemberAccountService struct {
	MemberID uuid.UUID `json:"member_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateMemberAccountService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_accounts.svc.update.start`)

	data, err := s.db.GetMemberAccountByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.MemberID != uuid.Nil {
		data.MemberID = req.MemberID
	}
	if req.Email != "" {
		data.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := hashing.HashPassword(req.Password)
		if err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}
		data.Password = string(hashedPassword)
	}

	if err := s.db.UpdateMemberAccount(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`member_accounts.svc.update.success`)
	return nil
}
