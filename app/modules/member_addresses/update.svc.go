package member_addresses

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdateMemberAddressService struct {
	MemberID      uuid.UUID `json:"member_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone"`
	IsDefault     *bool     `json:"is_default"`
	AddressNo     string    `json:"address_no"`
	Village       string    `json:"village"`
	Alley         string    `json:"alley"`
	SubDistrictID uuid.UUID `json:"sub_district_id"`
	DistrictID    uuid.UUID `json:"district_id"`
	ProvinceID    uuid.UUID `json:"province_id"`
	ZipcodeID     uuid.UUID `json:"zipcode_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateMemberAddressService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_addresses.svc.update.start`)

	data, err := s.db.GetMemberAddressByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.MemberID != uuid.Nil {
		data.MemberID = req.MemberID
	}
	if req.FirstName != "" {
		data.FirstName = req.FirstName
	}
	if req.LastName != "" {
		data.LastName = req.LastName
	}
	if req.Phone != "" {
		data.Phone = req.Phone
	}
	if req.IsDefault != nil {
		data.IsDefault = *req.IsDefault
	}
	if req.AddressNo != "" {
		data.AddressNo = req.AddressNo
	}
	if req.Village != "" {
		data.Village = req.Village
	}
	if req.Alley != "" {
		data.Alley = req.Alley
	}
	if req.SubDistrictID != uuid.Nil {
		data.SubDistrictID = req.SubDistrictID
	}
	if req.DistrictID != uuid.Nil {
		data.DistrictID = req.DistrictID
	}
	if req.ProvinceID != uuid.Nil {
		data.ProvinceID = req.ProvinceID
	}
	if req.ZipcodeID != uuid.Nil {
		data.ZipcodeID = req.ZipcodeID
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
			ActionType:   "member_address",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated member address " + data.ID.String(),
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

	span.AddEvent(`member_addresses.svc.update.success`)
	return nil
}
