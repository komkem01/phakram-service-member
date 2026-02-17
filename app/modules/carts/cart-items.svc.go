package carts

import (
	"context"
	"database/sql"
	"errors"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ListCartItemServiceRequest struct {
	base.RequestPaginate
	CartID uuid.UUID
}

type CreateCartItemServiceRequest struct {
	ProductID       uuid.UUID
	Quantity        int
	PricePerUnit    string
	TotalItemAmount string
}

type UpdateCartItemServiceRequest = CreateCartItemServiceRequest

func (s *Service) ListCartItemService(ctx context.Context, req *ListCartItemServiceRequest, requesterID uuid.UUID, isAdmin bool) ([]*ent.CartItemEntity, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.items.list.start`)

	if _, err := s.ensureCartAccess(ctx, req.CartID, requesterID, isAdmin); err != nil {
		return nil, nil, err
	}

	data := make([]*ent.CartItemEntity, 0)
	_, page, err := base.NewInstant(s.bunDB.DB()).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"cart_id", "product_id"},
		[]string{"created_at", "cart_id", "product_id"},
		func(selQ *bun.SelectQuery) *bun.SelectQuery {
			selQ.ExcludeColumn("price_per_unit", "total_item_amount")
			selQ.ColumnExpr("COALESCE(price_per_unit, 0) AS price_per_unit")
			selQ.ColumnExpr("COALESCE(total_item_amount, 0) AS total_item_amount")
			selQ.Where("cart_id = ?", req.CartID)
			return selQ
		},
	)
	if err != nil {
		return nil, nil, err
	}

	span.AddEvent(`carts.svc.items.list.success`)
	return data, page, nil
}

func (s *Service) InfoCartItemService(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID, requesterID uuid.UUID, isAdmin bool) (*ent.CartItemEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.items.info.start`)

	if _, err := s.ensureCartAccess(ctx, cartID, requesterID, isAdmin); err != nil {
		return nil, err
	}

	item, err := s.item.GetCartItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item.CartID != cartID {
		return nil, errors.New("cart items not found")
	}

	span.AddEvent(`carts.svc.items.info.success`)
	return item, nil
}

func (s *Service) CreateCartItemService(ctx context.Context, cartID uuid.UUID, req *CreateCartItemServiceRequest, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.items.create.start`)

	if _, err := s.ensureCartAccess(ctx, cartID, requesterID, isAdmin); err != nil {
		return err
	}

	if req.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	pricePerUnit, err := decimal.NewFromString(req.PricePerUnit)
	if err != nil {
		return err
	}
	totalItemAmount, err := parseTotalAmount(req.TotalItemAmount, pricePerUnit, req.Quantity)
	if err != nil {
		return err
	}

	existing := new(ent.CartItemEntity)
	err = s.bunDB.DB().NewSelect().
		Model(existing).
		ExcludeColumn("price_per_unit", "total_item_amount").
		ColumnExpr("COALESCE(price_per_unit, 0) AS price_per_unit").
		ColumnExpr("COALESCE(total_item_amount, 0) AS total_item_amount").
		Where("cart_id = ?", cartID).
		Where("product_id = ?", req.ProductID).
		Scan(ctx)
	if err == nil {
		existing.Quantity = req.Quantity
		existing.PricePerUnit = pricePerUnit
		existing.TotalItemAmount = totalItemAmount
		existing.UpdatedAt = time.Now()

		if err := s.item.UpdateCartItem(ctx, existing); err != nil {
			return err
		}

		span.AddEvent(`carts.svc.items.create.upserted`)
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	now := time.Now()
	item := &ent.CartItemEntity{
		ID:              uuid.New(),
		CartID:          cartID,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		PricePerUnit:    pricePerUnit,
		TotalItemAmount: totalItemAmount,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.item.CreateCartItem(ctx, item); err != nil {
		return err
	}

	span.AddEvent(`carts.svc.items.create.success`)
	return nil
}

func (s *Service) UpdateCartItemService(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID, req *UpdateCartItemServiceRequest, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.items.update.start`)

	if _, err := s.ensureCartAccess(ctx, cartID, requesterID, isAdmin); err != nil {
		return err
	}

	item, err := s.item.GetCartItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.CartID != cartID {
		return errors.New("cart items not found")
	}

	pricePerUnit, err := decimal.NewFromString(req.PricePerUnit)
	if err != nil {
		return err
	}
	totalItemAmount, err := parseTotalAmount(req.TotalItemAmount, pricePerUnit, req.Quantity)
	if err != nil {
		return err
	}

	item.ProductID = req.ProductID
	item.Quantity = req.Quantity
	item.PricePerUnit = pricePerUnit
	item.TotalItemAmount = totalItemAmount
	item.UpdatedAt = time.Now()

	if err := s.item.UpdateCartItem(ctx, item); err != nil {
		return err
	}

	span.AddEvent(`carts.svc.items.update.success`)
	return nil
}

func (s *Service) DeleteCartItemService(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.items.delete.start`)

	if _, err := s.ensureCartAccess(ctx, cartID, requesterID, isAdmin); err != nil {
		return err
	}

	item, err := s.item.GetCartItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.CartID != cartID {
		return errors.New("cart items not found")
	}

	if err := s.item.DeleteCartItem(ctx, itemID); err != nil {
		return err
	}

	span.AddEvent(`carts.svc.items.delete.success`)
	return nil
}

func parseTotalAmount(input string, pricePerUnit decimal.Decimal, quantity int) (decimal.Decimal, error) {
	if input == "" {
		return pricePerUnit.Mul(decimal.NewFromInt(int64(quantity))), nil
	}
	return decimal.NewFromString(input)
}
