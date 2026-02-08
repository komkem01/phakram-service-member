package members

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListMemberServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListMemberServiceResponses struct {
	ID            uuid.UUID       `json:"id"`
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
	Registration  *string         `json:"registration"`
	LastLogin     *string         `json:"last_login"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListMemberServiceRequest) ([]*ListMemberServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.list.start`)

	data, page, err := s.db.ListMembers(ctx, &entitiesdto.ListMembersRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListMemberServiceResponses
	for _, item := range data {
		temp := &ListMemberServiceResponses{
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
			Role:          string(item.Role),
			Phone:         item.Phone,
			TotalSpent:    item.TotalSpent,
			CurrentPoints: item.CurrentPoints,
			Registration: func() *string {
				if item.Registration != nil {
					str := item.Registration.Format("2006-01-02T15:04:05Z07:00")
					return &str
				}
				return nil
			}(),
			LastLogin: func() *string {
				if item.LastLogin != nil {
					str := item.LastLogin.Format("2006-01-02T15:04:05Z07:00")
					return &str
				}
				return nil
			}(),
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`members.svc.list.copy`)
	return response, page, nil
}
