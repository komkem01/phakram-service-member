package members

import (
	"context"
	"errors"
	"time"

	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateMemberWishlistServiceRequest struct {
	ProductID       uuid.UUID
	Quantity        int
	PricePerUnit    string
	TotalItemAmount string
	ActionBy        *uuid.UUID
}

type UpdateMemberWishlistServiceRequest = CreateMemberWishlistServiceRequest

func (s *Service) ListMemberWishlistService(ctx context.Context, req *entitiesdto.ListMemberWishlistRequest) ([]*ent.MemberWishlistEntity, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.wishlist.list.start`)

	data, page, err := s.wishlist.ListMemberWishlist(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	span.AddEvent(`members.svc.wishlist.list.success`)
	return data, page, nil
}

func (s *Service) CreateMemberWishlistService(ctx context.Context, memberID uuid.UUID, req *CreateMemberWishlistServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.wishlist.create.start`)

	pricePerUnit, err := decimal.NewFromString(req.PricePerUnit)
	if err != nil {
		return err
	}
	totalItemAmount, err := decimal.NewFromString(req.TotalItemAmount)
	if err != nil {
		return err
	}

	now := time.Now()
	wishlist := &ent.MemberWishlistEntity{
		ID:              uuid.New(),
		MemberID:        memberID,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		PricePerUnit:    pricePerUnit,
		TotalItemAmount: totalItemAmount,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(wishlist).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionCreated, ent.AuditActionCreated, "create_member_wishlist", wishlist.ID, req.ActionBy, "Created member wishlist with ID "+wishlist.ID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionCreated, "create_member_wishlist", wishlist.ID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.wishlist.create.success`)
	return nil
}

func (s *Service) InfoMemberWishlistService(ctx context.Context, memberID uuid.UUID, wishlistID uuid.UUID) (*ent.MemberWishlistEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.wishlist.info.start`)

	data, err := s.wishlist.GetMemberWishlistByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}
	if data.MemberID != memberID {
		return nil, errors.New("member wishlist not found")
	}

	span.AddEvent(`members.svc.wishlist.info.success`)
	return data, nil
}

func (s *Service) UpdateMemberWishlistService(ctx context.Context, memberID uuid.UUID, wishlistID uuid.UUID, req *UpdateMemberWishlistServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.wishlist.update.start`)

	pricePerUnit, err := decimal.NewFromString(req.PricePerUnit)
	if err != nil {
		return err
	}
	totalItemAmount, err := decimal.NewFromString(req.TotalItemAmount)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberWishlistEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", wishlistID).Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member wishlist not found")
		}

		data.ProductID = req.ProductID
		data.Quantity = req.Quantity
		data.PricePerUnit = pricePerUnit
		data.TotalItemAmount = totalItemAmount
		data.UpdatedAt = now
		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionUpdated, ent.AuditActionUpdated, "update_member_wishlist", data.ID, req.ActionBy, "Updated member wishlist with ID "+data.ID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionUpdated, "update_member_wishlist", wishlistID, req.ActionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.wishlist.update.success`)
	return nil
}

func (s *Service) DeleteMemberWishlistService(ctx context.Context, memberID uuid.UUID, wishlistID uuid.UUID, actionBy *uuid.UUID) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.wishlist.delete.start`)

	now := time.Now()
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.MemberWishlistEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", wishlistID).Scan(ctx); err != nil {
			return err
		}
		if data.MemberID != memberID {
			return errors.New("member wishlist not found")
		}

		if _, err := tx.NewDelete().Model(&ent.MemberWishlistEntity{}).Where("id = ?", wishlistID).Exec(ctx); err != nil {
			return err
		}
		return s.logMemberActionTx(ctx, tx, memberID, ent.MemberActionDeleted, ent.AuditActionDeleted, "delete_member_wishlist", wishlistID, actionBy, "Deleted member wishlist with ID "+wishlistID.String(), now)
	})
	if err != nil {
		s.logMemberActionFailed(ctx, ent.AuditActionDeleted, "delete_member_wishlist", wishlistID, actionBy, now, err)
		return err
	}

	span.AddEvent(`members.svc.wishlist.delete.success`)
	return nil
}
