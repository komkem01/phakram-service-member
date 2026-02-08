package members

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type UpdateMemberService struct {
	MemberNo      string           `json:"member_no"`
	TierID        uuid.UUID        `json:"tier_id"`
	StatusID      uuid.UUID        `json:"status_id"`
	PrefixID      uuid.UUID        `json:"prefix_id"`
	GenderID      uuid.UUID        `json:"gender_id"`
	FirstnameTh   string           `json:"firstname_th"`
	LastnameTh    string           `json:"lastname_th"`
	FirstnameEn   string           `json:"firstname_en"`
	LastnameEn    string           `json:"lastname_en"`
	Role          string           `json:"role"`
	Phone         string           `json:"phone"`
	TotalSpent    *decimal.Decimal `json:"total_spent"`
	CurrentPoints *int             `json:"current_points"`
	Registration  *time.Time       `json:"registration"`
	LastLogin     *time.Time       `json:"last_login"`
	MemberID      uuid.UUID        `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateMemberService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.update.start`)

	data, err := s.db.GetMemberByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.MemberNo != "" {
		data.MemberNo = req.MemberNo
	}
	if req.TierID != uuid.Nil {
		data.TierID = req.TierID
	}
	if req.StatusID != uuid.Nil {
		data.StatusID = req.StatusID
	}
	if req.PrefixID != uuid.Nil {
		data.PrefixID = req.PrefixID
	}
	if req.GenderID != uuid.Nil {
		data.GenderID = req.GenderID
	}
	if req.FirstnameTh != "" {
		data.FirstnameTh = req.FirstnameTh
	}
	if req.LastnameTh != "" {
		data.LastnameTh = req.LastnameTh
	}
	if req.FirstnameEn != "" {
		data.FirstnameEn = req.FirstnameEn
	}
	if req.LastnameEn != "" {
		data.LastnameEn = req.LastnameEn
	}
	if req.Role != "" {
		data.Role = ent.RoleTypeEnum(req.Role)
	}
	if req.Phone != "" {
		data.Phone = req.Phone
	}
	if req.TotalSpent != nil {
		data.TotalSpent = *req.TotalSpent
	}
	if req.CurrentPoints != nil {
		data.CurrentPoints = *req.CurrentPoints
	}
	if req.Registration != nil {
		data.Registration = req.Registration
	}
	if req.LastLogin != nil {
		data.LastLogin = req.LastLogin
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditUpdate,
			ActionType:   "member",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated member " + data.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`members.svc.update.success`)
	return nil
}
