package members

import (
	"context"
	"errors"
	entitiesdto "phakram/app/modules/entities/dto"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateMemberPaymentServiceRequest struct {
	PaymentID uuid.UUID
	Quantity  int
	Price     string
	ActionBy  *uuid.UUID
}

type UpdateMemberPaymentServiceRequest = CreateMemberPaymentServiceRequest

type ListMemberPaymentsServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

func (s *Service) ListMemberPaymentsService(ctx context.Context, req *ListMemberPaymentsServiceRequest) ([]*ent.MemberPaymentEntity, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.payment.list.start`)

	data, page, err := s.payment.ListMemberPayments(ctx, &entitiesdto.ListMemberPaymentsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, nil, err
	}

	span.AddEvent(`members.svc.payment.list.success`)
	return data, page, nil
}

func (s *Service) CreateMemberPaymentService(ctx context.Context, memberID uuid.UUID, req *CreateMemberPaymentServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.payment.create.start`)

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return err
	}

	now := time.Now()
	memberPayment := &ent.MemberPaymentEntity{
		ID:        uuid.New(),
		MemberID:  memberID,
		PaymentID: req.PaymentID,
		Quantity:  req.Quantity,
		Price:     price,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(memberPayment).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionCreated, ent.AuditActionCreated, "create_member_payment", memberPayment.ID, req.ActionBy, "Created member payment with ID "+memberPayment.ID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionCreated, "create_member_payment", memberPayment.ID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.payment.create.success`)
	return nil
}

func (s *Service) InfoMemberPaymentService(ctx context.Context, memberID uuid.UUID, rowID uuid.UUID) (*ent.MemberPaymentEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.payment.info.start`)

	data, err := s.payment.GetMemberPaymentByID(ctx, rowID)
	if err != nil {
		return nil, err
	}
	if data.MemberID != memberID {
		return nil, errors.New("member payment not found")
	}

	span.AddEvent(`members.svc.payment.info.success`)
	return data, nil
}

func (s *Service) UpdateMemberPaymentService(ctx context.Context, memberID uuid.UUID, rowID uuid.UUID, req *UpdateMemberPaymentServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.payment.update.start`)

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberPaymentEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", rowID).Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member payment not found")
		}

		data.PaymentID = req.PaymentID
		data.Quantity = req.Quantity
		data.Price = price
		data.UpdatedAt = now
		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionUpdated, ent.AuditActionUpdated, "update_member_payment", data.ID, req.ActionBy, "Updated member payment with ID "+data.ID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionUpdated, "update_member_payment", rowID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.payment.update.success`)
	return nil
}

func (s *Service) DeleteMemberPaymentService(ctx context.Context, memberID uuid.UUID, rowID uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.payment.delete.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberPaymentEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", rowID).Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member payment not found")
		}

		if _, err := tx.NewDelete().Model(&ent.MemberPaymentEntity{}).Where("id = ?", rowID).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionDeleted, ent.AuditActionDeleted, "delete_member_payment", rowID, actionBy, "Deleted member payment with ID "+rowID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionDeleted, "delete_member_payment", rowID, actionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.payment.delete.success`)
	return nil
}
