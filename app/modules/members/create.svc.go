package members

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateMemberService struct {
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
	TotalSpent    decimal.Decimal `json:"total_spent"`
	CurrentPoints int             `json:"current_points"`
	Registration  time.Time       `json:"registration"`
	LastLogin     time.Time       `json:"last_login"`
	MemberID      uuid.UUID       `json:"member_id"`
}

func (s *Service) CreateMemberService(ctx context.Context, req *CreateMemberService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.create.start`)

	memberNo := req.MemberNo
	if memberNo == "" {
		var err error
		memberNo, err = utils.GenerateMemberNo(ctx, s.bunDB.DB())
		if err != nil {
			return err
		}
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member := &ent.MemberEntity{
			ID:            uuid.New(),
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

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "member",
			ActionID:     &member.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created member " + member.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	span.AddEvent(`members.svc.create.success`)
	return nil
}
