package members

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListServiceRequest struct {
	base.RequestPaginate
	Role string
}

type ListServiceResponse struct {
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

func (s *Service) ListService(ctx context.Context, req *ListServiceRequest) ([]*ListServiceResponse, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.list.start`)

	data, page, err := s.db.ListMembers(ctx, &entitiesdto.ListMembersRequest{
		RequestPaginate: req.RequestPaginate,
		Role:            req.Role,
	})
	if err != nil {
		return nil, nil, err
	}

	resp := make([]*ListServiceResponse, 0, len(data))
	for _, item := range data {
		temp := &ListServiceResponse{
			ID:            item.ID,
			MemberNo:      item.MemberNo,
			TierID:        item.TierID,
			StatusID:      item.StatusID,
			PrefixID:      item.PrefixID,
			GenderID:      item.GenderID,
			FirstnameTh:   item.FirstnameTh,
			LastnameTh:    item.LastnameTh,
			FirstnameEn:   item.FirstnameEn,
			LastnameEn:    item.LastnameEn,
			Role:          item.Role,
			Phone:         item.Phone,
			TotalSpent:    item.TotalSpent,
			CurrentPoints: item.CurrentPoints,
			CreatedAt:     item.CreatedAt.Unix(),
			UpdatedAt:     item.UpdatedAt.Unix(),
		}

		if item.Registration != nil {
			registration := item.Registration.Unix()
			temp.Registration = &registration
		}
		if item.LastLogin != nil {
			lastLogin := item.LastLogin.Unix()
			temp.LastLogin = &lastLogin
		}

		resp = append(resp, temp)
	}

	span.AddEvent(`members.svc.list.success`)
	return resp, page, nil
}
