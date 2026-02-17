package orders

import (
	"context"
	"database/sql"
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
	MemberID           uuid.UUID
	PaymentID          uuid.UUID
	AddressID          uuid.UUID
	Status             string
	ShippingTrackingNo string
	TotalAmount        string
	DiscountAmount     string
	NetAmount          string
}

type UpdateOrderServiceRequest = CreateOrderServiceRequest

type ConfirmOrderPaymentServiceRequest struct {
	TransferredAmount string
	SlipImageBase64   string
	SlipFileName      string
	SlipFileType      string
	SlipFileSize      int64
}

type RejectOrderPaymentServiceRequest struct {
	Reason string
}

type OrderPaymentServiceResponse struct {
	OrderID       uuid.UUID `json:"order_id"`
	PaymentID     uuid.UUID `json:"payment_id"`
	OrderStatus   string    `json:"order_status"`
	PaymentStatus string    `json:"payment_status"`
	SlipAttached  bool      `json:"slip_attached"`
}

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

type orderPaymentReviewState struct {
	Submitted bool
	Rejected  bool
	Reason    string
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

	for _, item := range data {
		reviewState, reviewErr := s.getOrderPaymentReviewState(ctx, item.ID)
		if reviewErr != nil {
			return nil, nil, reviewErr
		}
		item.PaymentSubmitted = reviewState.Submitted
		item.PaymentRejected = reviewState.Rejected
		item.PaymentRejectionReason = reviewState.Reason

		trackingNo, trackingErr := s.getOrderShippingTrackingNo(ctx, item.ID)
		if trackingErr != nil {
			return nil, nil, trackingErr
		}
		item.ShippingTrackingNo = trackingNo
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

	paymentReviewState, err := s.getOrderPaymentReviewState(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	data.PaymentSubmitted = paymentReviewState.Submitted
	data.PaymentRejected = paymentReviewState.Rejected
	data.PaymentRejectionReason = paymentReviewState.Reason

	trackingNo, err := s.getOrderShippingTrackingNo(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	data.ShippingTrackingNo = trackingNo

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
		Where("action_type IN (?)", bun.In([]string{"order_status_transition", "order_payment_submitted", "order_payment_approved", "order_payment_rejected", "order_shipping_tracking_updated"})).
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

func (s *Service) CreateOrderService(ctx context.Context, req *CreateOrderServiceRequest) (*ent.OrderEntity, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.create.start`)

	totalAmount, err := decimal.NewFromString(req.TotalAmount)
	if err != nil {
		return nil, err
	}
	discountAmount, err := decimal.NewFromString(req.DiscountAmount)
	if err != nil {
		return nil, err
	}
	netAmount, err := decimal.NewFromString(req.NetAmount)
	if err != nil {
		return nil, err
	}

	requireMemberPayment := req.PaymentID != uuid.Nil
	if req.PaymentID == uuid.Nil {
		payment := &ent.PaymentEntity{
			ID:     uuid.New(),
			Amount: netAmount,
			Status: ent.PaymentTypePending,
		}
		if _, err := s.bunDB.DB().NewInsert().Model(payment).Exec(ctx); err != nil {
			return nil, err
		}
		req.PaymentID = payment.ID
	}

	if err := s.ensureCreateOrderOwnership(ctx, req.MemberID, req.AddressID, req.PaymentID, requireMemberPayment); err != nil {
		return nil, err
	}

	orderNo, err := utils.GenerateOrderNo()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	parsedStatus, err := parseOrderStatus(req.Status)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	span.AddEvent(`orders.svc.create.success`)
	return data, nil
}

func (s *Service) ensureCreateOrderOwnership(ctx context.Context, memberID uuid.UUID, addressID uuid.UUID, paymentID uuid.UUID, requireMemberPayment bool) error {
	address := new(ent.MemberAddressEntity)
	if err := s.bunDB.DB().NewSelect().Model(address).Where("id = ?", addressID).Scan(ctx); err != nil {
		return err
	}
	if address.MemberID != memberID {
		return errors.New("member address not found")
	}

	payment := new(ent.PaymentEntity)
	if err := s.bunDB.DB().NewSelect().Model(payment).Where("id = ?", paymentID).Scan(ctx); err != nil {
		return err
	}

	if !requireMemberPayment {
		return nil
	}

	memberPaymentCount, err := s.bunDB.DB().NewSelect().Model((*ent.MemberPaymentEntity)(nil)).Where("member_id = ?", memberID).Where("payment_id = ?", paymentID).Count(ctx)
	if err != nil {
		return err
	}
	if memberPaymentCount == 0 {
		return errors.New("member payment not found")
	}

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

	shippingTrackingNo := strings.TrimSpace(req.ShippingTrackingNo)

	paymentReviewState, err := s.getOrderPaymentReviewState(ctx, data.ID)
	if err != nil {
		return err
	}
	if data.Status == ent.StatusTypePending && paymentReviewState.Rejected {
		if nextStatus != ent.StatusTypePending {
			return errors.New("payment was rejected waiting for resubmission")
		}
	}

	if err := validateOrderStatusTransition(data.Status, nextStatus); err != nil {
		return err
	}

	previousStatus := data.Status
	statusChanged := previousStatus != nextStatus
	isShippingTransition := previousStatus != ent.StatusTypeShipping && nextStatus == ent.StatusTypeShipping
	if isShippingTransition && shippingTrackingNo == "" {
		return errors.New("shipping tracking number is required")
	}

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

			if isShippingTransition {
				trackingLog := &ent.AuditLogEntity{
					ID:           uuid.New(),
					Action:       ent.AuditActionUpdated,
					ActionType:   "order_shipping_tracking_updated",
					ActionID:     data.ID,
					ActionBy:     actionBy,
					Status:       ent.StatusAuditSuccesses,
					ActionDetail: "Shipping tracking number: " + shippingTrackingNo,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				if _, err := tx.NewInsert().Model(trackingLog).Exec(ctx); err != nil {
					return err
				}
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
		return []ent.StatusTypeEnum{ent.StatusTypeShipping}
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

func (s *Service) ConfirmOrderPaymentService(ctx context.Context, orderID uuid.UUID, req *ConfirmOrderPaymentServiceRequest, requesterID uuid.UUID, isAdmin bool) (*OrderPaymentServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.payment.confirm.start`)

	order, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin)
	if err != nil {
		return nil, err
	}

	if order.Status != ent.StatusTypePending {
		return nil, errors.New("order is not pending")
	}

	paymentReviewState, err := s.getOrderPaymentReviewState(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if paymentReviewState.Submitted {
		return nil, errors.New("payment confirmation already submitted")
	}

	paymentAmount := order.NetAmount
	if strings.TrimSpace(req.TransferredAmount) != "" {
		parsedAmount, err := decimal.NewFromString(strings.TrimSpace(req.TransferredAmount))
		if err != nil {
			return nil, err
		}
		paymentAmount = parsedAmount
	}

	slipAttached := strings.TrimSpace(req.SlipImageBase64) != ""
	uploadedBy := requesterID
	if uploadedBy == uuid.Nil {
		uploadedBy = order.MemberID
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		now := time.Now()

		paymentID := order.PaymentID
		if paymentID == uuid.Nil {
			paymentID = uuid.New()
			payment := &ent.PaymentEntity{
				ID:     paymentID,
				Amount: paymentAmount,
				Status: ent.PaymentTypePending,
			}
			if _, err := tx.NewInsert().Model(payment).Exec(ctx); err != nil {
				return err
			}

			order.PaymentID = paymentID
			order.UpdatedAt = now
			if _, err := tx.NewUpdate().Model(order).Where("id = ?", order.ID).Exec(ctx); err != nil {
				return err
			}
		} else {
			payment := new(ent.PaymentEntity)
			if err := tx.NewSelect().Model(payment).Where("id = ?", paymentID).Scan(ctx); err != nil {
				return err
			}

			payment.Amount = paymentAmount
			payment.Status = ent.PaymentTypePending
			payment.ApprovedBy = nil
			payment.ApprovedAt = nil

			if _, err := tx.NewUpdate().Model(payment).Where("id = ?", payment.ID).Exec(ctx); err != nil {
				return err
			}
		}

		if slipAttached {
			fileName := strings.TrimSpace(req.SlipFileName)
			if fileName == "" {
				fileName = fmt.Sprintf("payment-slip-%s", order.ID.String())
			}
			fileType := strings.TrimSpace(req.SlipFileType)
			if fileType == "" {
				fileType = "image/*"
			}

			storageID := uuid.New()
			storage := &ent.StorageEntity{
				ID:            storageID,
				RefID:         order.PaymentID,
				FileName:      fileName,
				FilePath:      strings.TrimSpace(req.SlipImageBase64),
				FileSize:      req.SlipFileSize,
				FileType:      fileType,
				IsActive:      true,
				RelatedEntity: ent.RelatedEntityPaymentFile,
				UploadedBy:    uploadedBy,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			if _, err := tx.NewInsert().Model(storage).Exec(ctx); err != nil {
				return err
			}

			paymentFile := &ent.PaymentFileEntity{
				ID:        uuid.New(),
				PaymentID: order.PaymentID,
				FileID:    storageID,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if _, err := tx.NewInsert().Model(paymentFile).Exec(ctx); err != nil {
				return err
			}
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "order_payment_submitted",
			ActionID:     order.ID,
			ActionBy:     &uploadedBy,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Payment confirmation submitted by member",
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`orders.svc.payment.confirm.success`)
	return &OrderPaymentServiceResponse{
		OrderID:       order.ID,
		PaymentID:     order.PaymentID,
		OrderStatus:   string(order.Status),
		PaymentStatus: string(ent.PaymentTypePending),
		SlipAttached:  slipAttached,
	}, nil
}

func (s *Service) getOrderPaymentReviewState(ctx context.Context, orderID uuid.UUID) (*orderPaymentReviewState, error) {
	state := &orderPaymentReviewState{}
	latestLog := new(ent.AuditLogEntity)

	err := s.bunDB.DB().NewSelect().
		Model(latestLog).
		Where("action_id = ?", orderID).
		Where("action_type IN (?)", bun.In([]string{"order_payment_submitted", "order_payment_approved", "order_payment_rejected"})).
		Where("status = ?", ent.StatusAuditSuccesses).
		OrderExpr("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return state, nil
		}
		return nil, err
	}

	switch latestLog.ActionType {
	case "order_payment_submitted":
		state.Submitted = true
	case "order_payment_rejected":
		state.Rejected = true
		state.Reason = parsePaymentRejectedReason(latestLog.ActionDetail)
	}

	return state, nil
}

func (s *Service) ApproveOrderPaymentService(ctx context.Context, orderID uuid.UUID, approverID uuid.UUID) (*OrderPaymentServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.payment.approve.start`)

	order, err := s.ensureOrderAccess(ctx, orderID, approverID, true)
	if err != nil {
		return nil, err
	}

	if order.PaymentID == uuid.Nil {
		return nil, errors.New("payment not found")
	}

	if order.Status != ent.StatusTypePending {
		return nil, errors.New("order is not pending")
	}

	paymentReviewState, err := s.getOrderPaymentReviewState(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if !paymentReviewState.Submitted {
		return nil, errors.New("payment confirmation not submitted")
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		now := time.Now()

		payment := new(ent.PaymentEntity)
		if err := tx.NewSelect().Model(payment).Where("id = ?", order.PaymentID).Scan(ctx); err != nil {
			return err
		}

		payment.Status = ent.PaymentTypeSuccess
		payment.ApprovedBy = &approverID
		payment.ApprovedAt = &now
		if _, err := tx.NewUpdate().Model(payment).Where("id = ?", payment.ID).Exec(ctx); err != nil {
			return err
		}

		previousStatus := order.Status
		order.Status = ent.StatusTypePaid
		order.UpdatedAt = now

		if err := s.applyOrderStatusSideEffects(ctx, tx, order, previousStatus, approverID); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().Model(order).Where("id = ?", order.ID).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "order_payment_approved",
			ActionID:     order.ID,
			ActionBy:     &approverID,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Order payment approved by admin",
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`orders.svc.payment.approve.success`)
	return &OrderPaymentServiceResponse{
		OrderID:       order.ID,
		PaymentID:     order.PaymentID,
		OrderStatus:   string(ent.StatusTypePaid),
		PaymentStatus: string(ent.PaymentTypeSuccess),
		SlipAttached:  false,
	}, nil
}

func (s *Service) RejectOrderPaymentService(ctx context.Context, orderID uuid.UUID, req *RejectOrderPaymentServiceRequest, approverID uuid.UUID) (*OrderPaymentServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.payment.reject.start`)

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return nil, errors.New("rejection reason is required")
	}

	order, err := s.ensureOrderAccess(ctx, orderID, approverID, true)
	if err != nil {
		return nil, err
	}

	if order.PaymentID == uuid.Nil {
		return nil, errors.New("payment not found")
	}

	if order.Status != ent.StatusTypePending {
		return nil, errors.New("order is not pending")
	}

	paymentReviewState, err := s.getOrderPaymentReviewState(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if !paymentReviewState.Submitted {
		return nil, errors.New("payment confirmation not submitted")
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		now := time.Now()

		payment := new(ent.PaymentEntity)
		if err := tx.NewSelect().Model(payment).Where("id = ?", order.PaymentID).Scan(ctx); err != nil {
			return err
		}

		payment.Status = ent.PaymentTypeFailed
		payment.ApprovedBy = nil
		payment.ApprovedAt = nil
		if _, err := tx.NewUpdate().Model(payment).Where("id = ?", payment.ID).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "order_payment_rejected",
			ActionID:     order.ID,
			ActionBy:     &approverID,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Payment rejected reason: " + reason,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`orders.svc.payment.reject.success`)
	return &OrderPaymentServiceResponse{
		OrderID:       order.ID,
		PaymentID:     order.PaymentID,
		OrderStatus:   string(order.Status),
		PaymentStatus: string(ent.PaymentTypeFailed),
		SlipAttached:  false,
	}, nil
}

func parsePaymentRejectedReason(detail string) string {
	const prefix = "Payment rejected reason: "
	if strings.HasPrefix(detail, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(detail, prefix))
	}
	return strings.TrimSpace(detail)
}

func (s *Service) getOrderShippingTrackingNo(ctx context.Context, orderID uuid.UUID) (string, error) {
	latestLog := new(ent.AuditLogEntity)
	err := s.bunDB.DB().NewSelect().
		Model(latestLog).
		Where("action_id = ?", orderID).
		Where("action_type = ?", "order_shipping_tracking_updated").
		Where("status = ?", ent.StatusAuditSuccesses).
		OrderExpr("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return parseShippingTrackingNumber(latestLog.ActionDetail), nil
}

func parseShippingTrackingNumber(detail string) string {
	const prefix = "Shipping tracking number: "
	if strings.HasPrefix(detail, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(detail, prefix))
	}
	return strings.TrimSpace(detail)
}
