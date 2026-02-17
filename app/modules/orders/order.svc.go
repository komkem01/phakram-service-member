package orders

import (
	"context"
	"errors"
	"fmt"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
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

type OrderTimelineItem struct {
	ActionType   string     `json:"action_type"`
	ActionDetail string     `json:"action_detail"`
	Status       string     `json:"status"`
	ActionBy     *uuid.UUID `json:"action_by"`
	FromStatus   string     `json:"from_status,omitempty"`
	ToStatus     string     `json:"to_status,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

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

func (s *Service) TimelineOrderService(ctx context.Context, orderID uuid.UUID, requesterID uuid.UUID, isAdmin bool) ([]*OrderTimelineItem, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.timeline.start`)

	if _, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin); err != nil {
		return nil, err
	}

	auditRows := make([]*ent.AuditLogEntity, 0)
	if err := s.bunDB.DB().NewSelect().
		Model(&auditRows).
		Where("action_id = ?", orderID).
		Where("action_type = ?", "order_status_transition").
		OrderExpr("created_at DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	items := make([]*OrderTimelineItem, 0, len(auditRows))
	for _, row := range auditRows {
		fromStatus, toStatus := parseOrderStatusTransitionDetail(row.ActionDetail)
		items = append(items, &OrderTimelineItem{
			ActionType:   row.ActionType,
			ActionDetail: row.ActionDetail,
			Status:       string(row.Status),
			ActionBy:     row.ActionBy,
			FromStatus:   fromStatus,
			ToStatus:     toStatus,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		})
	}

	span.AddEvent(`orders.svc.timeline.success`)
	return items, nil
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
	parsedStatus, err := parseOrderStatus(req.Status)
	if err != nil {
		return err
	}

	data := &ent.OrderEntity{
		ID:             uuid.New(),
		OrderNo:        orderNo,
		MemberID:       req.MemberID,
		PaymentID:      req.PaymentID,
		AddressID:      req.AddressID,
		Status:         parsedStatus,
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

	if strings.TrimSpace(req.Status) == "" {
		return errors.New("status is required")
	}

	nextStatus, err := parseOrderStatus(req.Status)
	if err != nil {
		return err
	}

	if err := validateOrderStatusTransition(data.Status, nextStatus); err != nil {
		return err
	}

	previousStatus := data.Status
	statusChanged := previousStatus != nextStatus

	data.PaymentID = req.PaymentID
	data.AddressID = req.AddressID
	data.Status = nextStatus
	data.TotalAmount = totalAmount
	data.DiscountAmount = discountAmount
	data.NetAmount = netAmount
	data.UpdatedAt = time.Now()

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if statusChanged {
			if err := s.applyOrderStatusSideEffects(ctx, tx, data, previousStatus, requesterID); err != nil {
				return err
			}
		}

		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		if statusChanged {
			auditDetail := fmt.Sprintf("Order status changed from %s to %s", previousStatus, data.Status)
			var actionBy *uuid.UUID
			if requesterID != uuid.Nil {
				actionBy = &requesterID
			}

			auditLog := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.AuditActionUpdated,
				ActionType:   "order_status_transition",
				ActionID:     data.ID,
				ActionBy:     actionBy,
				Status:       ent.StatusAuditSuccesses,
				ActionDetail: auditDetail,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
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

func parseOrderStatus(status string) (ent.StatusTypeEnum, error) {
	value := strings.ToLower(strings.TrimSpace(status))
	if value == "" {
		return ent.StatusTypePending, nil
	}

	switch value {
	case string(ent.StatusTypePaid):
		return ent.StatusTypePaid, nil
	case string(ent.StatusTypeShipping):
		return ent.StatusTypeShipping, nil
	case string(ent.StatusTypeCompleted):
		return ent.StatusTypeCompleted, nil
	case string(ent.StatusTypeCancelled):
		return ent.StatusTypeCancelled, nil
	case string(ent.StatusTypePending):
		return ent.StatusTypePending, nil
	default:
		return "", errors.New("invalid order status")
	}
}

func validateOrderStatusTransition(currentStatus ent.StatusTypeEnum, nextStatus ent.StatusTypeEnum) error {
	if currentStatus == nextStatus {
		return nil
	}

	allowedStatuses := allowedNextStatuses(currentStatus)
	for _, allowedStatus := range allowedStatuses {
		if allowedStatus == nextStatus {
			return nil
		}
	}

	return errors.New("invalid status transition")
}

func allowedNextStatuses(status ent.StatusTypeEnum) []ent.StatusTypeEnum {
	switch status {
	case ent.StatusTypePending:
		return []ent.StatusTypeEnum{ent.StatusTypePaid, ent.StatusTypeCancelled}
	case ent.StatusTypePaid:
		return []ent.StatusTypeEnum{ent.StatusTypeShipping, ent.StatusTypeCancelled}
	case ent.StatusTypeShipping:
		return []ent.StatusTypeEnum{ent.StatusTypeCompleted}
	default:
		return []ent.StatusTypeEnum{}
	}
}

func parseOrderStatusTransitionDetail(detail string) (string, string) {
	const prefix = "Order status changed from "
	if !strings.HasPrefix(detail, prefix) {
		return "", ""
	}

	parts := strings.SplitN(strings.TrimPrefix(detail, prefix), " to ", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func (s *Service) applyOrderStatusSideEffects(ctx context.Context, tx bun.Tx, order *ent.OrderEntity, previousStatus ent.StatusTypeEnum, requesterID uuid.UUID) error {
	if previousStatus != ent.StatusTypePaid && order.Status == ent.StatusTypePaid {
		if err := s.createMemberPaymentFromOrder(ctx, tx, order); err != nil {
			return err
		}
	}

	if previousStatus != ent.StatusTypeShipping && order.Status == ent.StatusTypeShipping {
		if err := s.decreaseStockFromOrderItems(ctx, tx, order.ID); err != nil {
			return err
		}
	}

	_ = requesterID
	return nil
}

func (s *Service) createMemberPaymentFromOrder(ctx context.Context, tx bun.Tx, order *ent.OrderEntity) error {
	items, err := s.listOrderItemsByOrderID(ctx, tx, order.ID)
	if err != nil {
		return err
	}

	totalQuantity := 0
	for _, item := range items {
		totalQuantity += item.Quantity
	}

	payment := &ent.MemberPaymentEntity{
		ID:        uuid.New(),
		MemberID:  order.MemberID,
		PaymentID: order.PaymentID,
		Quantity:  totalQuantity,
		Price:     order.NetAmount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := tx.NewInsert().Model(payment).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) decreaseStockFromOrderItems(ctx context.Context, tx bun.Tx, orderID uuid.UUID) error {
	items, err := s.listOrderItemsByOrderID(ctx, tx, orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		stock := new(ent.ProductStockEntity)
		if err := tx.NewSelect().
			Model(stock).
			Where("product_id = ?", item.ProductID).
			Where("deleted_at IS NULL").
			Limit(1).
			For("UPDATE").
			Scan(ctx); err != nil {
			return err
		}

		if stock.Remaining < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s", item.ProductID.String())
		}

		stock.Remaining -= item.Quantity
		stock.UpdatedAt = time.Now()

		if _, err := tx.NewUpdate().Model(stock).Where("id = ?", stock.ID).Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) listOrderItemsByOrderID(ctx context.Context, db bun.IDB, orderID uuid.UUID) ([]*ent.OrderItemEntity, error) {
	items := make([]*ent.OrderItemEntity, 0)
	if err := db.NewSelect().Model(&items).Where("order_id = ?", orderID).Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
