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

type CreateMemberBankServiceRequest struct {
	BankID      uuid.UUID
	BankNo      string
	FirstnameTh string
	LastnameTh  string
	FirstnameEn string
	LastnameEn  string
	IsDefault   bool
	ActionBy    *uuid.UUID
}

type UpdateMemberBankServiceRequest = CreateMemberBankServiceRequest

func (s *Service) CreateMemberBankService(ctx context.Context, memberID uuid.UUID, req *CreateMemberBankServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.bank.create.start`)

	now := time.Now()
	memberBank := &ent.MemberBankEntity{
		ID:          uuid.New(),
		MemberID:    memberID,
		BankID:      req.BankID,
		BankNo:      req.BankNo,
		FirstnameTh: req.FirstnameTh,
		LastnameTh:  req.LastnameTh,
		FirstnameEn: req.FirstnameEn,
		LastnameEn:  req.LastnameEn,
		IsDefault:   req.IsDefault,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(memberBank).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(
			ctx,
			tx,
			memberID,
			ent.MemberActionCreated,
			ent.AuditActionCreated,
			"create_member_bank",
			memberBank.ID,
			req.ActionBy,
			"Created member bank with ID "+memberBank.ID.String(),
			now,
		)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionCreated, "create_member_bank", memberBank.ID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.bank.create.success`)
	return nil
}

func (s *Service) InfoMemberBankService(ctx context.Context, memberID uuid.UUID, memberBankID uuid.UUID) (*ent.MemberBankEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.bank.info.start`)

	data, err := s.bank.GetMemberBankByID(ctx, memberBankID)
	if err != nil {
		return nil, err
	}
	if data.MemberID != memberID {
		return nil, errors.New("member bank not found")
	}

	span.AddEvent(`members.svc.bank.info.success`)
	return data, nil
}

func (s *Service) UpdateMemberBankService(ctx context.Context, memberID uuid.UUID, memberBankID uuid.UUID, req *UpdateMemberBankServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.bank.update.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberBankEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", memberBankID).Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member bank not found")
		}

		data.BankID = req.BankID
		data.BankNo = req.BankNo
		data.FirstnameTh = req.FirstnameTh
		data.LastnameTh = req.LastnameTh
		data.FirstnameEn = req.FirstnameEn
		data.LastnameEn = req.LastnameEn
		data.IsDefault = req.IsDefault
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
			"update_member_bank",
			data.ID,
			req.ActionBy,
			"Updated member bank with ID "+data.ID.String(),
			now,
		)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionUpdated, "update_member_bank", memberBankID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.bank.update.success`)
	return nil
}

func (s *Service) DeleteMemberBankService(ctx context.Context, memberID uuid.UUID, memberBankID uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.bank.delete.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberBankEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", memberBankID).Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member bank not found")
		}

		if _, err := tx.NewDelete().Model(&ent.MemberBankEntity{}).Where("id = ?", memberBankID).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(
			ctx,
			tx,
			memberID,
			ent.MemberActionDeleted,
			ent.AuditActionDeleted,
			"delete_member_bank",
			memberBankID,
			actionBy,
			"Deleted member bank with ID "+memberBankID.String(),
			now,
		)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionDeleted, "delete_member_bank", memberBankID, actionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.bank.delete.success`)
	return nil
}
