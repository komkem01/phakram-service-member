package member_addresses

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
)

type InfoMemberAddressServiceResponses struct {
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

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoMemberAddressServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_addresses.svc.info.start`)

	data, err := s.db.GetMemberAddressByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoMemberAddressServiceResponses{
		ID:            data.ID,
		MemberID:      data.MemberID,
		FirstName:     data.FirstName,
		LastName:      data.LastName,
		Phone:         data.Phone,
		IsDefault:     data.IsDefault,
		AddressNo:     data.AddressNo,
		Village:       data.Village,
		Alley:         data.Alley,
		SubDistrictID: data.SubDistrictID,
		DistrictID:    data.DistrictID,
		ProvinceID:    data.ProvinceID,
		ZipcodeID:     data.ZipcodeID,
		CreatedAt:     data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_addresses.svc.info.success`)
	return resp, nil
}
