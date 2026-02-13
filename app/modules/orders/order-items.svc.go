package orders

import (
	"context"
	"errors"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type ListOrderItemServiceRequest struct {
	base.RequestPaginate
	OrderID uuid.UUID
}

type CreateOrderItemServiceRequest struct {
	ProductID       uuid.UUID
	Quantity        int
	PricePerUnit    string
	TotalItemAmount string
}

type UpdateOrderItemServiceRequest = CreateOrderItemServiceRequest

func (s *Service) ListOrderItemService(ctx context.Context, req *ListOrderItemServiceRequest, requesterID uuid.UUID, isAdmin bool) ([]*ent.OrderItemEntity, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.items.list.start`)

	if _, err := s.ensureOrderAccess(ctx, req.OrderID, requesterID, isAdmin); err != nil {
		return nil, nil, err
	}

	data := make([]*ent.OrderItemEntity, 0)
	_, page, err := base.NewInstant(s.bunDB.DB()).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"order_id", "product_id"},
		[]string{"created_at", "order_id", "product_id"},
		func(selQ *bun.SelectQuery) *bun.SelectQuery {
			selQ.Where("order_id = ?", req.OrderID)
			return selQ
		},
	)
	if err != nil {
		return nil, nil, err
	}

	span.AddEvent(`orders.svc.items.list.success`)
	return data, page, nil
}

func (s *Service) InfoOrderItemService(ctx context.Context, orderID uuid.UUID, itemID uuid.UUID, requesterID uuid.UUID, isAdmin bool) (*ent.OrderItemEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.items.info.start`)

	if _, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin); err != nil {
		return nil, err
	}

	item, err := s.item.GetOrderItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item.OrderID != orderID {
		return nil, errors.New("order item not found")
	}

	span.AddEvent(`orders.svc.items.info.success`)
	return item, nil
}

func (s *Service) CreateOrderItemService(ctx context.Context, orderID uuid.UUID, req *CreateOrderItemServiceRequest, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.items.create.start`)

	if _, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin); err != nil {
		return err
	}

	pricePerUnit, err := decimal.NewFromString(req.PricePerUnit)
	if err != nil {
		return err
	}
	totalItemAmount, err := parseOrderItemTotalAmount(req.TotalItemAmount, pricePerUnit, req.Quantity)
	if err != nil {
		return err
	}

	now := time.Now()
	item := &ent.OrderItemEntity{
		ID:              uuid.New(),
		OrderID:         orderID,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		PricePerUnit:    pricePerUnit,
		TotalItemAmount: totalItemAmount,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.item.CreateOrderItem(ctx, item); err != nil {
		return err
	}

	span.AddEvent(`orders.svc.items.create.success`)
	return nil
}

func (s *Service) UpdateOrderItemService(ctx context.Context, orderID uuid.UUID, itemID uuid.UUID, req *UpdateOrderItemServiceRequest, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.items.update.start`)

	if _, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin); err != nil {
		return err
	}

	item, err := s.item.GetOrderItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.OrderID != orderID {
		return errors.New("order item not found")
	}

	pricePerUnit, err := decimal.NewFromString(req.PricePerUnit)
	if err != nil {
		return err
	}
	totalItemAmount, err := parseOrderItemTotalAmount(req.TotalItemAmount, pricePerUnit, req.Quantity)
	if err != nil {
		return err
	}

	item.ProductID = req.ProductID
	item.Quantity = req.Quantity
	item.PricePerUnit = pricePerUnit
	item.TotalItemAmount = totalItemAmount
	item.UpdatedAt = time.Now()

	if err := s.item.UpdateOrderItem(ctx, item); err != nil {
		return err
	}

	span.AddEvent(`orders.svc.items.update.success`)
	return nil
}

func (s *Service) DeleteOrderItemService(ctx context.Context, orderID uuid.UUID, itemID uuid.UUID, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.items.delete.start`)

	if _, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin); err != nil {
		return err
	}

	item, err := s.item.GetOrderItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.OrderID != orderID {
		return errors.New("order item not found")
	}

	if err := s.item.DeleteOrderItem(ctx, itemID); err != nil {
		return err
	}

	span.AddEvent(`orders.svc.items.delete.success`)
	return nil
}

func parseOrderItemTotalAmount(input string, pricePerUnit decimal.Decimal, quantity int) (decimal.Decimal, error) {
	if input == "" {
		return pricePerUnit.Mul(decimal.NewFromInt(int64(quantity))), nil
	}
	return decimal.NewFromString(input)
}
