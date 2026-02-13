package members

import (
	"context"
	"errors"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateMemberAddressServiceRequest struct {
	FirstName     string
	LastName      string
	Phone         string
	IsDefault     bool
	AddressNo     string
	Village       string
	Alley         string
	SubDistrictID uuid.UUID
	DistrictID    uuid.UUID
	ProvinceID    uuid.UUID
	ZipcodeID     uuid.UUID
	ActionBy      *uuid.UUID
}

type UpdateMemberAddressServiceRequest = CreateMemberAddressServiceRequest

func (s *Service) CreateMemberAddressService(ctx context.Context, memberID uuid.UUID, req *CreateMemberAddressServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.address.create.start`)

	now := time.Now()
	address := &ent.MemberAddressEntity{
		ID:            uuid.New(),
		MemberID:      memberID,
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
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(address).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(
			ctx,
			tx,
			memberID,
			ent.MemberActionCreated,
			ent.AuditActionCreated,
			"create_member_address",
			address.ID,
			req.ActionBy,
			"Created member address with ID "+address.ID.String(),
			now,
		)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionCreated, "create_member_address", address.ID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.address.create.success`)
	return nil
}

func (s *Service) InfoMemberAddressService(ctx context.Context, memberID uuid.UUID, addressID uuid.UUID) (*ent.MemberAddressEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.address.info.start`)

	data, err := s.address.GetMemberAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if data.MemberID != memberID {
		return nil, errors.New("member address not found")
	}

	span.AddEvent(`members.svc.address.info.success`)
	return data, nil
}

func (s *Service) UpdateMemberAddressService(ctx context.Context, memberID uuid.UUID, addressID uuid.UUID, req *UpdateMemberAddressServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.address.update.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberAddressEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", addressID).Where("deleted_at IS NULL").Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member address not found")
		}

		data.FirstName = req.FirstName
		data.LastName = req.LastName
		data.Phone = req.Phone
		data.IsDefault = req.IsDefault
		data.AddressNo = req.AddressNo
		data.Village = req.Village
		data.Alley = req.Alley
		data.SubDistrictID = req.SubDistrictID
		data.DistrictID = req.DistrictID
		data.ProvinceID = req.ProvinceID
		data.ZipcodeID = req.ZipcodeID
		data.UpdatedAt = now

		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(
			ctx,
			tx,
			memberID,
			ent.MemberActionUpdated,
			ent.AuditActionUpdated,
			"update_member_address",
			data.ID,
			req.ActionBy,
			"Updated member address with ID "+data.ID.String(),
			now,
		)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionUpdated, "update_member_address", addressID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.address.update.success`)
	return nil
}

func (s *Service) DeleteMemberAddressService(ctx context.Context, memberID uuid.UUID, addressID uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.address.delete.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberAddressEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", addressID).Where("deleted_at IS NULL").Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member address not found")
		}

		if _, err := tx.NewUpdate().Model(&ent.MemberAddressEntity{}).Set("deleted_at = ?", now).Set("updated_at = ?", now).Where("id = ?", addressID).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(
			ctx,
			tx,
			memberID,
			ent.MemberActionDeleted,
			ent.AuditActionDeleted,
			"delete_member_address",
			addressID,
			actionBy,
			"Deleted member address with ID "+addressID.String(),
			now,
		)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionDeleted, "delete_member_address", addressID, actionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.address.delete.success`)
	return nil
}
