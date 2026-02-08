package member_addresses

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateMemberAddressService struct {
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
}

func (s *Service) CreateMemberAddressService(ctx context.Context, req *CreateMemberAddressService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_addresses.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		address := &ent.MemberAddressEntity{
			ID:            uuid.New(),
			MemberID:      req.MemberID,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			Phone:         req.Phone,
			IsDefault:     req.IsDefault,
			AddressNo:     req.AddressNo,
			Village:       req.Village,
			Alley:         req.Alley,
			SubDistrictID: req.SubDistrictID,
			DistrictID:    req.DistrictID,
			ProvinceID:    req.ProvinceID,
			ZipcodeID:     req.ZipcodeID,
		}
		if _, err := tx.NewInsert().Model(address).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "member_address",
			ActionID:     &address.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created member address " + address.ID.String(),
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
	span.AddEvent(`member_addresses.svc.create.success`)
	return nil
}
