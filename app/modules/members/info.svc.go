package members

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoServiceResponse struct {
	ID            uuid.UUID        `json:"id"`
	MemberNo      string           `json:"member_no"`
	TierID        uuid.UUID        `json:"tier_id"`
	StatusID      uuid.UUID        `json:"status_id"`
	PrefixID      uuid.UUID        `json:"prefix_id"`
	GenderID      uuid.UUID        `json:"gender_id"`
	FirstnameTh   string           `json:"firstname_th"`
	LastnameTh    string           `json:"lastname_th"`
	FirstnameEn   string           `json:"firstname_en"`
	LastnameEn    string           `json:"lastname_en"`
	Role          ent.RoleTypeEnum `json:"role"`
	Phone         string           `json:"phone"`
	TotalSpent    decimal.Decimal  `json:"total_spent"`
	CurrentPoints int              `json:"current_points"`
	CreatedAt     int64            `json:"created_at"`
	UpdatedAt     int64            `json:"updated_at"`
	Registration  *int64           `json:"registration"`
	LastLogin     *int64           `json:"last_login"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.info.start`)

	data, err := s.db.GetMemberByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &InfoServiceResponse{
		ID:            data.ID,
		MemberNo:      data.MemberNo,
		TierID:        data.TierID,
		StatusID:      data.StatusID,
		PrefixID:      data.PrefixID,
		GenderID:      data.GenderID,
		FirstnameTh:   data.FirstnameTh,
		LastnameTh:    data.LastnameTh,
		FirstnameEn:   data.FirstnameEn,
		LastnameEn:    data.LastnameEn,
		Role:          data.Role,
		Phone:         data.Phone,
		TotalSpent:    data.TotalSpent,
		CurrentPoints: data.CurrentPoints,
		CreatedAt:     data.CreatedAt.Unix(),
		UpdatedAt:     data.UpdatedAt.Unix(),
	}
	if data.Registration != nil {
		registration := data.Registration.Unix()
		resp.Registration = &registration
	}
	if data.LastLogin != nil {
		lastLogin := data.LastLogin.Unix()
		resp.LastLogin = &lastLogin
	}

	span.AddEvent(`members.svc.info.success`)
	return resp, nil
}
