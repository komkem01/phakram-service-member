package orders

import (
	"context"
	"errors"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListOrderServiceRequest struct {
	base.RequestPaginate
	MemberID  uuid.UUID
	Search    string
	Status    string
	StartDate int64
	EndDate   int64
}

type CreateOrderServiceRequest struct {
	MemberID       uuid.UUID
	PaymentID      uuid.UUID
	AddressID      uuid.UUID
	Status         string
	TotalAmount    string
	DiscountAmount string
	NetAmount      string
}

type UpdateOrderServiceRequest = CreateOrderServiceRequest

func (s *Service) ListOrderService(ctx context.Context, req *ListOrderServiceRequest) ([]*ent.OrderEntity, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.list.start`)

	data, page, err := s.order.ListOrders(ctx, &entitiesdto.ListOrdersRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
		Search:          req.Search,
		Status:          req.Status,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
	})
	if err != nil {
		return nil, nil, err
	}

	span.AddEvent(`orders.svc.list.success`)
	return data, page, nil
}

func (s *Service) InfoOrderService(ctx context.Context, orderID uuid.UUID, requesterID uuid.UUID, isAdmin bool) (*ent.OrderEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.info.start`)

	data, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin)
	if err != nil {
		return nil, err
	}

	span.AddEvent(`orders.svc.info.success`)
	return data, nil
}

func (s *Service) CreateOrderService(ctx context.Context, req *CreateOrderServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.create.start`)

	totalAmount, err := decimal.NewFromString(req.TotalAmount)
	if err != nil {
		return err
	}
	discountAmount, err := decimal.NewFromString(req.DiscountAmount)
	if err != nil {
		return err
	}
	netAmount, err := decimal.NewFromString(req.NetAmount)
	if err != nil {
		return err
	}

	orderNo, err := utils.GenerateOrderNo()
	if err != nil {
		return err
	}

	now := time.Now()
	data := &ent.OrderEntity{
		ID:             uuid.New(),
		OrderNo:        orderNo,
		MemberID:       req.MemberID,
		PaymentID:      req.PaymentID,
		AddressID:      req.AddressID,
		Status:         parseOrderStatus(req.Status),
		TotalAmount:    totalAmount,
		DiscountAmount: discountAmount,
		NetAmount:      netAmount,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.order.CreateOrder(ctx, data); err != nil {
		return err
	}

	span.AddEvent(`orders.svc.create.success`)
	return nil
}

func (s *Service) UpdateOrderService(ctx context.Context, orderID uuid.UUID, req *UpdateOrderServiceRequest, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.update.start`)

	data, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin)
	if err != nil {
		return err
	}

	totalAmount, err := decimal.NewFromString(req.TotalAmount)
	if err != nil {
		return err
	}
	discountAmount, err := decimal.NewFromString(req.DiscountAmount)
	if err != nil {
		return err
	}
	netAmount, err := decimal.NewFromString(req.NetAmount)
	if err != nil {
		return err
	}

	data.PaymentID = req.PaymentID
	data.AddressID = req.AddressID
	data.Status = parseOrderStatus(req.Status)
	data.TotalAmount = totalAmount
	data.DiscountAmount = discountAmount
	data.NetAmount = netAmount
	data.UpdatedAt = time.Now()

	if err := s.order.UpdateOrder(ctx, data); err != nil {
		return err
	}

	span.AddEvent(`orders.svc.update.success`)
	return nil
}

func (s *Service) DeleteOrderService(ctx context.Context, orderID uuid.UUID, requesterID uuid.UUID, isAdmin bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.delete.start`)

	if _, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin); err != nil {
		return err
	}

	if err := s.order.DeleteOrder(ctx, orderID); err != nil {
		return err
	}

	span.AddEvent(`orders.svc.delete.success`)
	return nil
}

func (s *Service) ensureOrderAccess(ctx context.Context, orderID uuid.UUID, requesterID uuid.UUID, isAdmin bool) (*ent.OrderEntity, error) {
	data, err := s.order.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if !isAdmin && data.MemberID != requesterID {
		return nil, errors.New("forbidden")
	}
	if data.ID == uuid.Nil {
		return nil, errors.New("order not found")
	}
	return data, nil
}

func parseOrderStatus(status string) ent.StatusTypeEnum {
	value := strings.ToLower(strings.TrimSpace(status))
	switch value {
	case string(ent.StatusTypePaid):
		return ent.StatusTypePaid
	case string(ent.StatusTypeShipping):
		return ent.StatusTypeShipping
	case string(ent.StatusTypeCompleted):
		return ent.StatusTypeCompleted
	case string(ent.StatusTypeCancelled):
		return ent.StatusTypeCancelled
	default:
		return ent.StatusTypePending
	}
}
