package member_addresses

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberAddressServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListMemberAddressServiceResponses struct {
	ID            uuid.UUID `json:"id"`
	MemberID      uuid.UUID `json:"member_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone"`
	IsDefault     bool      `json:"is_default"`
	AddressNo     string    `json:"address_no"`
	Village       string    `json:"village"`
	Alley         string    `json:"alley"`
	SubDistrictID uuid.UUID `json:"sub_district_id"`
	DistrictID    uuid.UUID `json:"district_id"`
	ProvinceID    uuid.UUID `json:"province_id"`
	ZipcodeID     uuid.UUID `json:"zipcode_id"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListMemberAddressServiceRequest) ([]*ListMemberAddressServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_addresses.svc.list.start`)

	data, page, err := s.db.ListMemberAddresses(ctx, &entitiesdto.ListMemberAddressesRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListMemberAddressServiceResponses
	for _, item := range data {
		temp := &ListMemberAddressServiceResponses{
			ID:            item.ID,
			MemberID:      item.MemberID,
			FirstName:     item.FirstName,
			LastName:      item.LastName,
			Phone:         item.Phone,
			IsDefault:     item.IsDefault,
			AddressNo:     item.AddressNo,
			Village:       item.Village,
			Alley:         item.Alley,
			SubDistrictID: item.SubDistrictID,
			DistrictID:    item.DistrictID,
			ProvinceID:    item.ProvinceID,
			ZipcodeID:     item.ZipcodeID,
			CreatedAt:     item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`member_addresses.svc.list.copy`)
	return response, page, nil
}
