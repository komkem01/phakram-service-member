package carts

import (
	"context"
	"errors"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"time"

	"github.com/google/uuid"
)

type ListCartServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type CreateCartServiceRequest struct {
	MemberID uuid.UUID
	IsActive *bool
}

type UpdateCartServiceRequest struct {
	IsActive *bool
}

func (s *Service) ListCartService(ctx context.Context, req *ListCartServiceRequest) ([]*ent.CartEntity, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.list.start`)

	data, page, err := s.cart.ListCarts(ctx, &entitiesdto.ListCartsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, nil, err
	}

	span.AddEvent(`carts.svc.list.success`)
	return data, page, nil
}

func (s *Service) InfoCartService(ctx context.Context, cartID uuid.UUID, requesterID uuid.UUID, isAdmin bool) (*ent.CartEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.info.start`)

	data, err := s.ensureCartAccess(ctx, cartID, requesterID, isAdmin)
	if err != nil {
		return nil, err
	}

	span.AddEvent(`carts.svc.info.success`)
	return data, nil
}

func (s *Service) CreateCartService(ctx context.Context, req *CreateCartServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.create.start`)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	now := time.Now()
	data := &ent.CartEntity{
		ID:        uuid.New(),
		MemberID:  req.MemberID,
		IsActive:  isActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.cart.CreateCart(ctx, data); err != nil {
		return err
	}

	span.AddEvent(`carts.svc.create.success`)
	return nil
}

func (s *Service) UpdateCartService(ctx context.Context, cartID uuid.UUID, req *UpdateCartServiceRequest, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.update.start`)

	data, err := s.ensureCartAccess(ctx, cartID, requesterID, isAdmin)
	if err != nil {
		return err
	}

	if req.IsActive != nil {
		data.IsActive = *req.IsActive
	}
	data.UpdatedAt = time.Now()

	if err := s.cart.UpdateCart(ctx, data); err != nil {
		return err
	}

	span.AddEvent(`carts.svc.update.success`)
	return nil
}

func (s *Service) DeleteCartService(ctx context.Context, cartID uuid.UUID, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.delete.start`)

	if _, err := s.ensureCartAccess(ctx, cartID, requesterID, isAdmin); err != nil {
		return err
	}

	if err := s.cart.DeleteCart(ctx, cartID); err != nil {
		return err
	}

	span.AddEvent(`carts.svc.delete.success`)
	return nil
}

func (s *Service) ensureCartAccess(ctx context.Context, cartID uuid.UUID, requesterID uuid.UUID, isAdmin bool) (*ent.CartEntity, error) {
	data, err := s.cart.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if !isAdmin && data.MemberID != requesterID {
		return nil, errors.New("forbidden")
	}
	if data.ID == uuid.Nil {
		return nil, errors.New("cart not found")
	}
	return data, nil
}
