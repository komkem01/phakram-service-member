package orders

import (
	"context"
	"database/sql"
	"encoding/base64"
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

type ListMemberNotificationServiceRequest struct {
	base.RequestPaginate
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

type UpdateOrderServiceRequest struct {
	PaymentID          uuid.UUID
	AddressID          uuid.UUID
	Status             string
	ShippingTrackingNo string
	TotalAmount        string
	DiscountAmount     string
	NetAmount          string
	CancelReason       string
	RefundReason       string
}

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

type AppealOrderPaymentServiceRequest struct {
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
	ActionType         string     `json:"action_type"`
	ActionDetail       string     `json:"action_detail"`
	Status             string     `json:"status"`
	ActionBy           *uuid.UUID `json:"action_by"`
	FromStatus         string     `json:"from_status,omitempty"`
	ToStatus           string     `json:"to_status,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type MemberNotificationItem struct {
	ID          uuid.UUID `json:"id"`
	EventType   string    `json:"event_type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	OrderID     uuid.UUID `json:"order_id"`
	OrderNo     string    `json:"order_no"`
	OrderStatus string    `json:"order_status"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type MemberNotificationMarkAllReadResponse struct {
	UpdatedCount int `json:"updated_count"`
}

type MemberNotificationUnreadCountResponse struct {
	UnreadCount int `json:"unread_count"`
}

type ReorderServiceResponse struct {
	OrderID           uuid.UUID `json:"order_id"`
	CartID            uuid.UUID `json:"cart_id"`
	AddedProductCount int       `json:"added_product_count"`
}

type orderPaymentReviewState struct {
	Submitted bool
	Rejected  bool
	Reason    string
}

const (
	orderPaymentReviewStatusSubmitted = "submitted"
	orderPaymentReviewStatusApproved  = "approved"
	orderPaymentReviewStatusRejected  = "rejected"
)

type memberNotificationRow struct {
	ID           uuid.UUID          `bun:"id"`
	ActionType   string             `bun:"action_type"`
	ActionDetail string             `bun:"action_detail"`
	CreatedAt    time.Time          `bun:"created_at"`
	OrderID      uuid.UUID          `bun:"order_id"`
	OrderNo      string             `bun:"order_no"`
	OrderStatus  ent.StatusTypeEnum `bun:"order_status"`
	IsRead       bool               `bun:"is_read"`
}

type memberNotificationReadEntity struct {
	bun.BaseModel  `bun:"table:member_notification_reads,alias:mnr"`
	MemberID       uuid.UUID `bun:"member_id,pk,type:uuid"`
	NotificationID uuid.UUID `bun:"notification_id,pk,type:uuid"`
	ReadAt         time.Time `bun:"read_at"`
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
		item.PaymentRejectionReason = normalizePaymentRejectionReason(reviewState.Reason)

		trackingNo, trackingErr := s.getOrderShippingTrackingNo(ctx, item.ID)
		if trackingErr != nil {
			return nil, nil, trackingErr
		}
		item.ShippingTrackingNo = trackingNo
		item.StatusSummary, item.StatusNextStep = mapOrderStatusSummary(item.Status, reviewState.Submitted, reviewState.Rejected)
	}

	span.AddEvent(`orders.svc.list.success`)
	return data, page, nil
}

func (s *Service) ListMemberNotificationService(ctx context.Context, req *ListMemberNotificationServiceRequest, requesterID uuid.UUID, isAdmin bool) ([]*MemberNotificationItem, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.notifications.list.start`)

	if requesterID == uuid.Nil {
		return nil, nil, errors.New("forbidden")
	}

	allowedEventTypes := memberNotificationEventTypes()

	query := s.bunDB.DB().NewSelect().
		TableExpr("audit_log AS al").
		Join("JOIN orders AS o ON o.id = al.action_id").
		Join("LEFT JOIN member_notification_reads AS mnr ON mnr.notification_id = al.id AND mnr.member_id = ?", requesterID).
		Where("al.status = ?", ent.StatusAuditSuccesses).
		Where("al.action_type IN (?)", bun.In(allowedEventTypes))

	_ = isAdmin
	query = query.Where("o.member_id = ?", requesterID)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	rows := make([]*memberNotificationRow, 0)
	offset := (req.GetPage() - 1) * req.GetSize()
	if err := query.
		ColumnExpr("al.id AS id").
		ColumnExpr("al.action_type AS action_type").
		ColumnExpr("al.action_detail AS action_detail").
		ColumnExpr("al.created_at AS created_at").
		ColumnExpr("o.id AS order_id").
		ColumnExpr("o.order_no AS order_no").
		ColumnExpr("o.status AS order_status").
		ColumnExpr("CASE WHEN mnr.member_id IS NULL THEN FALSE ELSE TRUE END AS is_read").
		OrderExpr("al.created_at DESC").
		Offset(int(offset)).
		Limit(int(req.GetSize())).
		Scan(ctx, &rows); err != nil {
		return nil, nil, err
	}

	items := make([]*MemberNotificationItem, 0, len(rows))
	for _, row := range rows {
		title, message := mapNotificationTitleMessage(row.ActionType, row.ActionDetail, row.OrderNo)
		items = append(items, &MemberNotificationItem{
			ID:          row.ID,
			EventType:   row.ActionType,
			Title:       title,
			Message:     message,
			OrderID:     row.OrderID,
			OrderNo:     row.OrderNo,
			OrderStatus: string(row.OrderStatus),
			IsRead:      row.IsRead,
			CreatedAt:   row.CreatedAt,
		})
	}

	page := &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}
	span.AddEvent(`orders.svc.notifications.list.success`)
	return items, page, nil
}

func (s *Service) MarkMemberNotificationReadService(ctx context.Context, notificationID uuid.UUID, requesterID uuid.UUID, isAdmin bool, isRead bool) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.notifications.mark.start`)

	if err := s.ensureMemberNotificationAccess(ctx, notificationID, requesterID, isAdmin); err != nil {
		return err
	}

	if isRead {
		entry := &memberNotificationReadEntity{
			MemberID:       requesterID,
			NotificationID: notificationID,
			ReadAt:         time.Now(),
		}

		if _, err := s.bunDB.DB().NewInsert().
			Model(entry).
			On("CONFLICT (member_id, notification_id) DO UPDATE").
			Set("read_at = EXCLUDED.read_at").
			Exec(ctx); err != nil {
			return err
		}
	} else {
		if _, err := s.bunDB.DB().NewDelete().
			Model((*memberNotificationReadEntity)(nil)).
			Where("member_id = ?", requesterID).
			Where("notification_id = ?", notificationID).
			Exec(ctx); err != nil {
			return err
		}
	}

	span.AddEvent(`orders.svc.notifications.mark.success`)
	return nil
}

func (s *Service) MarkAllMemberNotificationsReadService(ctx context.Context, requesterID uuid.UUID, isAdmin bool) (*MemberNotificationMarkAllReadResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.notifications.mark_all.start`)

	if requesterID == uuid.Nil {
		return nil, errors.New("forbidden")
	}

	allowedEventTypes := memberNotificationEventTypes()

	notificationIDs := make([]uuid.UUID, 0)
	query := s.bunDB.DB().NewSelect().
		TableExpr("audit_log AS al").
		Join("JOIN orders AS o ON o.id = al.action_id").
		Where("al.status = ?", ent.StatusAuditSuccesses).
		Where("al.action_type IN (?)", bun.In(allowedEventTypes))

	_ = isAdmin
	query = query.Where("o.member_id = ?", requesterID)

	if err := query.
		ColumnExpr("al.id").
		Scan(ctx, &notificationIDs); err != nil {
		return nil, err
	}

	if len(notificationIDs) == 0 {
		return &MemberNotificationMarkAllReadResponse{UpdatedCount: 0}, nil
	}

	now := time.Now()
	entries := make([]*memberNotificationReadEntity, 0, len(notificationIDs))
	for _, notificationID := range notificationIDs {
		entries = append(entries, &memberNotificationReadEntity{
			MemberID:       requesterID,
			NotificationID: notificationID,
			ReadAt:         now,
		})
	}

	if _, err := s.bunDB.DB().NewInsert().
		Model(&entries).
		On("CONFLICT (member_id, notification_id) DO UPDATE").
		Set("read_at = EXCLUDED.read_at").
		Exec(ctx); err != nil {
		return nil, err
	}

	span.AddEvent(`orders.svc.notifications.mark_all.success`)
	return &MemberNotificationMarkAllReadResponse{UpdatedCount: len(entries)}, nil
}

func (s *Service) CountMemberUnreadNotificationsService(ctx context.Context, requesterID uuid.UUID, isAdmin bool) (*MemberNotificationUnreadCountResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.notifications.count_unread.start`)

	if requesterID == uuid.Nil {
		return nil, errors.New("forbidden")
	}

	allowedEventTypes := memberNotificationEventTypes()

	count, err := s.bunDB.DB().NewSelect().
		TableExpr("audit_log AS al").
		Join("JOIN orders AS o ON o.id = al.action_id").
		Join("LEFT JOIN member_notification_reads AS mnr ON mnr.notification_id = al.id AND mnr.member_id = ?", requesterID).
		Where("al.status = ?", ent.StatusAuditSuccesses).
		Where("al.action_type IN (?)", bun.In(allowedEventTypes)).
		Where("o.member_id = ?", requesterID).
		Where("mnr.notification_id IS NULL").
		Count(ctx)
	if err != nil {
		return nil, err
	}

	_ = isAdmin

	span.AddEvent(`orders.svc.notifications.count_unread.success`)
	return &MemberNotificationUnreadCountResponse{UnreadCount: count}, nil
}

func (s *Service) ensureMemberNotificationAccess(ctx context.Context, notificationID uuid.UUID, requesterID uuid.UUID, isAdmin bool) error {
	if requesterID == uuid.Nil {
		return errors.New("forbidden")
	}

	allowedEventTypes := memberNotificationEventTypes()

	query := s.bunDB.DB().NewSelect().
		TableExpr("audit_log AS al").
		Join("JOIN orders AS o ON o.id = al.action_id").
		Where("al.id = ?", notificationID).
		Where("al.status = ?", ent.StatusAuditSuccesses).
		Where("al.action_type IN (?)", bun.In(allowedEventTypes))

	_ = isAdmin
	query = query.Where("o.member_id = ?", requesterID)

	count, err := query.Count(ctx)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("notification not found")
	}

	return nil
}

func memberNotificationEventTypes() []string {
	return []string{
		"order_payment_appealed",
		"order_payment_approved",
		"order_payment_rejected",
		"order_refund_rejected",
		"order_status_transition",
		"order_shipping_tracking_updated",
	}
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
	data.PaymentRejectionReason = normalizePaymentRejectionReason(paymentReviewState.Reason)

	trackingNo, err := s.getOrderShippingTrackingNo(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	data.ShippingTrackingNo = trackingNo
	data.StatusSummary, data.StatusNextStep = mapOrderStatusSummary(data.Status, paymentReviewState.Submitted, paymentReviewState.Rejected)

	cancellationReason, err := s.getOrderCancellationReason(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	data.CancellationReason = cancellationReason

	refundRejectionReason, err := s.getOrderRefundRejectionReason(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	data.RefundRejectionReason = refundRejectionReason

	paymentAppealReason, err := s.getOrderPaymentAppealReason(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	data.PaymentAppealReason = paymentAppealReason

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
		Where("action_type IN (?)", bun.In([]string{"order_status_transition", "order_payment_submitted", "order_payment_appealed", "order_payment_approved", "order_payment_rejected", "order_refund_rejected", "order_shipping_tracking_updated"})).
		OrderExpr("created_at DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	items := make([]*OrderTimelineItem, 0, len(auditRows))
	cancellationReason := ""
	cancellationReasonLoaded := false
	for _, row := range auditRows {
		fromStatus, toStatus := parseOrderStatusTransitionDetail(row.ActionDetail)
		itemCancellationReason := ""
		if row.ActionType == "order_status_transition" && toStatus == string(ent.StatusTypeCancelled) {
			if !cancellationReasonLoaded {
				loadedReason, reasonErr := s.getOrderCancellationReason(ctx, orderID)
				if reasonErr != nil {
					return nil, reasonErr
				}
				cancellationReason = loadedReason
				cancellationReasonLoaded = true
			}
			itemCancellationReason = cancellationReason
		}

		items = append(items, &OrderTimelineItem{
			ActionType:         row.ActionType,
			ActionDetail:       row.ActionDetail,
			Status:             string(row.Status),
			ActionBy:           row.ActionBy,
			FromStatus:         fromStatus,
			ToStatus:           toStatus,
			CancellationReason: itemCancellationReason,
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt,
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
	discountAmount, netAmount, err := s.calculateOrderAmountsByMemberTier(ctx, req.MemberID, totalAmount)
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

func (s *Service) calculateOrderAmountsByMemberTier(ctx context.Context, memberID uuid.UUID, totalAmount decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	discountRate, err := s.getMemberTierDiscountRate(ctx, memberID)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	if discountRate.IsNegative() {
		discountRate = decimal.Zero
	}
	maxRate := decimal.NewFromInt(100)
	if discountRate.GreaterThan(maxRate) {
		discountRate = maxRate
	}

	discountAmount := totalAmount.Mul(discountRate).Div(maxRate).Round(2)
	if discountAmount.GreaterThan(totalAmount) {
		discountAmount = totalAmount
	}

	netAmount := totalAmount.Sub(discountAmount).Round(2)
	if netAmount.IsNegative() {
		netAmount = decimal.Zero
	}

	return discountAmount, netAmount, nil
}

func (s *Service) getMemberTierDiscountRate(ctx context.Context, memberID uuid.UUID) (decimal.Decimal, error) {
	member := new(ent.MemberEntity)
	if err := s.bunDB.DB().NewSelect().
		Model(member).
		Column("tier_id").
		Where("id = ?", memberID).
		Limit(1).
		Scan(ctx); err != nil {
		return decimal.Zero, err
	}

	if member.TierID == uuid.Nil {
		return decimal.Zero, nil
	}

	tier := new(ent.TierEntity)
	err := s.bunDB.DB().NewSelect().
		Model(tier).
		Column("discount_rate").
		Where("id = ?", member.TierID).
		Where("is_active = ?", true).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return decimal.Zero, nil
		}
		return decimal.Zero, err
	}

	return tier.DiscountRate, nil
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
	cancelReason := normalizeOrderCancellationReason(req.CancelReason)
	refundReason := normalizeOrderCancellationReason(req.RefundReason)
	isRequestingRefund := nextStatus == ent.StatusTypeRefundRequested
	isRefundReviewTransition := data.Status == ent.StatusTypeRefundRequested && (nextStatus == ent.StatusTypePaid || nextStatus == ent.StatusTypeCancelled)

	paymentReviewState, err := s.getOrderPaymentReviewState(ctx, data.ID)
	if err != nil {
		return err
	}

	if isRequestingRefund {
		if !isAdmin {
			if data.Status != ent.StatusTypePaid && data.Status != ent.StatusTypeShipping && data.Status != ent.StatusTypeCompleted {
				return errors.New("refund request is allowed only after payment approval")
			}
		}

		paymentStatus, paymentErr := s.getOrderPaymentStatus(ctx, data.PaymentID)
		if paymentErr != nil {
			return paymentErr
		}
		if paymentStatus != ent.PaymentTypeSuccess {
			return errors.New("refund request requires successful payment")
		}
	}

	if isRefundReviewTransition {
		if !isAdmin {
			return errors.New("only admin can review refund request")
		}

		paymentStatus, paymentErr := s.getOrderPaymentStatus(ctx, data.PaymentID)
		if paymentErr != nil {
			return paymentErr
		}
		if paymentStatus != ent.PaymentTypeSuccess {
			return errors.New("refund review requires successful payment")
		}
	}

	if nextStatus == ent.StatusTypeCancelled && paymentReviewState.Submitted && data.Status != ent.StatusTypeRefundRequested {
		return errors.New("cannot cancel order after payment submission")
	}
	if nextStatus == ent.StatusTypeRefundRequested && !paymentReviewState.Submitted {
		return errors.New("refund request requires payment submission")
	}
	if nextStatus == ent.StatusTypeRefundRequested && paymentReviewState.Rejected {
		return errors.New("payment was rejected waiting for resubmission")
	}
	if nextStatus == ent.StatusTypeRefundRequested && refundReason == "" {
		return errors.New("refund reason is required")
	}
	if data.Status == ent.StatusTypeRefundRequested && nextStatus == ent.StatusTypePaid && refundReason == "" {
		return errors.New("refund rejection reason is required")
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
	isRefundRejectionTransition := previousStatus == ent.StatusTypeRefundRequested && nextStatus == ent.StatusTypePaid
	isShippingTransition := previousStatus != ent.StatusTypeShipping && nextStatus == ent.StatusTypeShipping
	if isShippingTransition && shippingTrackingNo == "" {
		return errors.New("shipping tracking number is required")
	}
	if nextStatus == ent.StatusTypeCancelled && cancelReason == "" {
		existingReason, reasonErr := s.getOrderCancellationReason(ctx, data.ID)
		if reasonErr != nil {
			return reasonErr
		}
		cancelReason = existingReason
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
			if nextStatus == ent.StatusTypeRefundRequested {
				if err := s.upsertOrderCancellationInTx(ctx, tx, data.ID, requesterID, isAdmin, refundReason); err != nil {
					return err
				}
			}

			if nextStatus == ent.StatusTypeCancelled {
				if err := s.upsertOrderCancellationInTx(ctx, tx, data.ID, requesterID, isAdmin, cancelReason); err != nil {
					return err
				}
			}

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

				if err := s.upsertOrderShippingTrackingInTx(ctx, tx, data.ID, shippingTrackingNo, actionBy); err != nil {
					return err
				}
			}

			if isRefundRejectionTransition {
				refundRejectedLog := &ent.AuditLogEntity{
					ID:           uuid.New(),
					Action:       ent.AuditActionUpdated,
					ActionType:   "order_refund_rejected",
					ActionID:     data.ID,
					ActionBy:     actionBy,
					Status:       ent.StatusAuditSuccesses,
					ActionDetail: "Refund rejected reason: " + refundReason,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				if _, err := tx.NewInsert().Model(refundRejectedLog).Exec(ctx); err != nil {
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

func (s *Service) upsertOrderCancellationInTx(ctx context.Context, tx bun.Tx, orderID uuid.UUID, requesterID uuid.UUID, isAdmin bool, reason string) error {
	now := time.Now()
	var cancelledBy *uuid.UUID
	if requesterID != uuid.Nil {
		cancelledBy = &requesterID
	}

	cancelledRole := "member"
	if isAdmin {
		cancelledRole = "admin"
	}

	record := &ent.OrderCancellationEntity{
		ID:            uuid.New(),
		OrderID:       orderID,
		CancelledBy:   cancelledBy,
		CancelledRole: cancelledRole,
		Reason:        reason,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if _, err := tx.NewInsert().
		Model(record).
		On("CONFLICT (order_id) DO UPDATE").
		Set("cancelled_by = EXCLUDED.cancelled_by").
		Set("cancelled_role = EXCLUDED.cancelled_role").
		Set("reason = EXCLUDED.reason").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func normalizeOrderCancellationReason(reason string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(reason)), " ")
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

func (s *Service) ReorderService(ctx context.Context, orderID uuid.UUID, requesterID uuid.UUID, isAdmin bool) (*ReorderServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.reorder.start`)

	order, err := s.ensureOrderAccess(ctx, orderID, requesterID, isAdmin)
	if err != nil {
		return nil, err
	}
	if order.Status != ent.StatusTypeCompleted {
		return nil, errors.New("reorder is allowed only for completed orders")
	}

	orderItems := make([]*ent.OrderItemEntity, 0)
	if err := s.bunDB.DB().NewSelect().Model(&orderItems).Where("order_id = ?", order.ID).Scan(ctx); err != nil {
		return nil, err
	}
	if len(orderItems) == 0 {
		return nil, errors.New("order has no items")
	}

	result := &ReorderServiceResponse{OrderID: order.ID}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		cartID, err := s.ensureActiveCartIDInTx(ctx, tx, order.MemberID)
		if err != nil {
			return err
		}
		result.CartID = cartID

		for _, orderItem := range orderItems {
			product := new(ent.ProductEntity)
			if err := tx.NewSelect().Model(product).Where("id = ?", orderItem.ProductID).Scan(ctx); err != nil {
				return err
			}
			if !product.IsActive {
				return errors.New("product is inactive")
			}

			stock := new(ent.ProductStockEntity)
			if err := tx.NewSelect().
				Model(stock).
				Where("product_id = ?", orderItem.ProductID).
				Where("deleted_at IS NULL").
				OrderExpr("updated_at DESC").
				Limit(1).
				Scan(ctx); err != nil {
				return err
			}

			existingCartItem := new(ent.CartItemEntity)
			err := tx.NewSelect().
				Model(existingCartItem).
				Where("cart_id = ?", cartID).
				Where("product_id = ?", orderItem.ProductID).
				Limit(1).
				Scan(ctx)

			now := time.Now()
			if errors.Is(err, sql.ErrNoRows) {
				if stock.Remaining < orderItem.Quantity {
					return errors.New("insufficient product stock")
				}

				cartItem := &ent.CartItemEntity{
					ID:              uuid.New(),
					CartID:          cartID,
					ProductID:       orderItem.ProductID,
					Quantity:        orderItem.Quantity,
					PricePerUnit:    orderItem.PricePerUnit,
					TotalItemAmount: orderItem.PricePerUnit.Mul(decimal.NewFromInt(int64(orderItem.Quantity))),
					CreatedAt:       now,
					UpdatedAt:       now,
				}
				if _, err := tx.NewInsert().Model(cartItem).Exec(ctx); err != nil {
					return err
				}
				result.AddedProductCount++
				continue
			}
			if err != nil {
				return err
			}

			nextQuantity := existingCartItem.Quantity + orderItem.Quantity
			if stock.Remaining < nextQuantity {
				return errors.New("insufficient product stock")
			}

			existingCartItem.Quantity = nextQuantity
			existingCartItem.TotalItemAmount = existingCartItem.PricePerUnit.Mul(decimal.NewFromInt(int64(nextQuantity)))
			existingCartItem.UpdatedAt = now

			if _, err := tx.NewUpdate().Model(existingCartItem).Where("id = ?", existingCartItem.ID).Exec(ctx); err != nil {
				return err
			}
			result.AddedProductCount++
		}

		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`orders.svc.reorder.success`)
	return result, nil
}

func (s *Service) ensureActiveCartIDInTx(ctx context.Context, tx bun.Tx, memberID uuid.UUID) (uuid.UUID, error) {
	activeCart := new(ent.CartEntity)
	err := tx.NewSelect().
		Model(activeCart).
		Where("member_id = ?", memberID).
		Where("is_active = ?", true).
		OrderExpr("updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err == nil {
		return activeCart.ID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, err
	}

	now := time.Now()
	createdCart := &ent.CartEntity{
		ID:        uuid.New(),
		MemberID:  memberID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := tx.NewInsert().Model(createdCart).Exec(ctx); err != nil {
		return uuid.Nil, err
	}

	return createdCart.ID, nil
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
	case string(ent.StatusTypeRefundRequested):
		return ent.StatusTypeRefundRequested, nil
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
		return []ent.StatusTypeEnum{ent.StatusTypePaid, ent.StatusTypeCancelled, ent.StatusTypeRefundRequested}
	case ent.StatusTypeRefundRequested:
		return []ent.StatusTypeEnum{ent.StatusTypePaid, ent.StatusTypeCancelled}
	case ent.StatusTypePaid:
		return []ent.StatusTypeEnum{ent.StatusTypeShipping, ent.StatusTypeRefundRequested}
	case ent.StatusTypeShipping:
		return []ent.StatusTypeEnum{ent.StatusTypeCompleted, ent.StatusTypeRefundRequested}
	case ent.StatusTypeCompleted:
		return []ent.StatusTypeEnum{ent.StatusTypeRefundRequested}
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
	if previousStatus == ent.StatusTypePending && order.Status == ent.StatusTypePaid {
		if err := s.createMemberPaymentFromOrder(ctx, tx, order); err != nil {
			return err
		}
	}

	if previousStatus == ent.StatusTypeRefundRequested && order.Status == ent.StatusTypeCancelled {
		if order.PaymentID == uuid.Nil {
			return errors.New("payment not found")
		}

		payment := new(ent.PaymentEntity)
		if err := tx.NewSelect().Model(payment).Where("id = ?", order.PaymentID).Limit(1).Scan(ctx); err != nil {
			return err
		}

		now := time.Now()
		payment.Status = ent.PaymentTypeRefunded
		if requesterID != uuid.Nil {
			payment.ApprovedBy = &requesterID
			payment.ApprovedAt = &now
		} else {
			payment.ApprovedBy = nil
			payment.ApprovedAt = nil
		}

		if _, err := tx.NewUpdate().Model(payment).Where("id = ?", payment.ID).Exec(ctx); err != nil {
			return err
		}
	}

	if previousStatus != ent.StatusTypeShipping && order.Status == ent.StatusTypeShipping {
		if err := s.decreaseStockFromOrderItems(ctx, tx, order.ID); err != nil {
			return err
		}
	}

	if previousStatus != ent.StatusTypeCompleted && order.Status == ent.StatusTypeCompleted {
		if err := s.addMemberSpendAndPointsFromOrder(ctx, tx, order); err != nil {
			return err
		}
	}

	_ = requesterID
	return nil
}

func (s *Service) addMemberSpendAndPointsFromOrder(ctx context.Context, tx bun.Tx, order *ent.OrderEntity) error {
	actualPaidAmount, err := s.getActualPaidAmountInTx(ctx, tx, order)
	if err != nil {
		return err
	}

	member := new(ent.MemberEntity)
	if err := tx.NewSelect().
		Model(member).
		Where("id = ?", order.MemberID).
		For("UPDATE").
		Limit(1).
		Scan(ctx); err != nil {
		return err
	}

	earnedPoints := int(actualPaidAmount.Div(decimal.NewFromInt(100)).Floor().IntPart())
	if earnedPoints < 0 {
		earnedPoints = 0
	}

	now := time.Now()
	member.TotalSpent = member.TotalSpent.Add(actualPaidAmount).Round(2)
	member.CurrentPoints += earnedPoints

	upgradedTierName := ""
	upgradedToTierID := uuid.Nil
	nextTier, err := s.findHighestEligibleTierBySpendingInTx(ctx, tx, member.TotalSpent)
	if err != nil {
		return err
	}
	if nextTier != nil && member.TierID != nextTier.ID {
		member.TierID = nextTier.ID
		upgradedToTierID = nextTier.ID
		if strings.TrimSpace(nextTier.NameTh) != "" {
			upgradedTierName = strings.TrimSpace(nextTier.NameTh)
		} else {
			upgradedTierName = strings.TrimSpace(nextTier.NameEn)
		}
	}

	member.UpdatedAt = now

	if _, err := tx.NewUpdate().
		Model(member).
		Column("total_spent", "current_points", "tier_id", "updated_at").
		Where("id = ?", member.ID).
		Exec(ctx); err != nil {
		return err
	}

	details := fmt.Sprintf("Order %s completed: total_spent +%s, points +%d", order.OrderNo, actualPaidAmount.StringFixed(2), earnedPoints)
	if upgradedToTierID != uuid.Nil {
		if upgradedTierName != "" {
			details = fmt.Sprintf("%s, tier upgraded to %s", details, upgradedTierName)
		} else {
			details = fmt.Sprintf("%s, tier upgraded", details)
		}
	}

	memberTx := &ent.MemberTransactionEntity{
		ID:        uuid.New(),
		MemberID:  member.ID,
		Action:    ent.MemberActionUpdated,
		Details:   details,
		CreatedAt: now,
	}
	if _, err := tx.NewInsert().Model(memberTx).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) findHighestEligibleTierBySpendingInTx(ctx context.Context, tx bun.Tx, totalSpent decimal.Decimal) (*ent.TierEntity, error) {
	tier := new(ent.TierEntity)
	err := tx.NewSelect().
		Model(tier).
		Where("is_active = ?", true).
		Where("min_spending <= ?", totalSpent).
		OrderExpr("min_spending DESC, created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return tier, nil
}

func (s *Service) getActualPaidAmountInTx(ctx context.Context, tx bun.Tx, order *ent.OrderEntity) (decimal.Decimal, error) {
	if order.PaymentID == uuid.Nil {
		return order.NetAmount.Round(2), nil
	}

	payment := new(ent.PaymentEntity)
	err := tx.NewSelect().
		Model(payment).
		Where("id = ?", order.PaymentID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order.NetAmount.Round(2), nil
		}
		return decimal.Zero, err
	}

	amount := payment.Amount.Round(2)
	if amount.IsNegative() {
		return decimal.Zero, nil
	}

	return amount, nil
}

func (s *Service) getOrderPaymentStatus(ctx context.Context, paymentID uuid.UUID) (ent.PaymentTypeEnum, error) {
	if paymentID == uuid.Nil {
		return "", errors.New("payment not found")
	}

	payment := new(ent.PaymentEntity)
	if err := s.bunDB.DB().NewSelect().Model(payment).Where("id = ?", paymentID).Limit(1).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("payment not found")
		}
		return "", err
	}

	return payment.Status, nil
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

	if _, err := tx.NewInsert().
		Model(payment).
		On("CONFLICT (member_id, payment_id) DO UPDATE").
		Set("quantity = EXCLUDED.quantity").
		Set("price = EXCLUDED.price").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx); err != nil {
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

	slipFilePath := ""
	slipFileName := strings.TrimSpace(req.SlipFileName)
	slipFileType := strings.TrimSpace(req.SlipFileType)
	slipFileSize := req.SlipFileSize

	if slipAttached {
		trimmedSlipBase64 := strings.TrimSpace(req.SlipImageBase64)

		if s.supabase != nil && s.supabase.enabledForPrivate() {
			uploadedSlip, err := s.supabase.UploadPaymentSlip(ctx, order.ID, order.PaymentID, slipFileName, trimmedSlipBase64)
			if err != nil {
				return nil, err
			}

			slipFilePath = uploadedSlip.Path
			if slipFileName == "" {
				slipFileName = uploadedSlip.FileName
			}
			if slipFileType == "" {
				slipFileType = uploadedSlip.MIMEType
			}
			if slipFileSize <= 0 {
				slipFileSize = uploadedSlip.Size
			}
		} else {
			decodedSlip, detectedMIME, err := decodeBase64Image(trimmedSlipBase64)
			if err != nil {
				return nil, err
			}
			if len(decodedSlip) == 0 {
				return nil, errors.New("slip image is empty")
			}
			if len(decodedSlip) > maxSlipFileSizeBytes {
				return nil, errors.New("slip image exceeds 5 MB")
			}
			if !isAllowedImageMIME(detectedMIME) {
				return nil, fmt.Errorf("unsupported slip image type: %s", detectedMIME)
			}

			slipFilePath = trimmedSlipBase64
			if !strings.HasPrefix(strings.ToLower(trimmedSlipBase64), "data:") {
				slipFilePath = fmt.Sprintf("data:%s;base64,%s", detectedMIME, base64.StdEncoding.EncodeToString(decodedSlip))
			}

			if slipFileType == "" {
				slipFileType = detectedMIME
			}
			if slipFileSize <= 0 {
				slipFileSize = int64(len(decodedSlip))
			}
		}
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
			fileName := strings.TrimSpace(slipFileName)
			if fileName == "" {
				fileName = fmt.Sprintf("payment-slip-%s", order.ID.String())
			}
			fileType := strings.TrimSpace(slipFileType)
			if fileType == "" {
				fileType = "image/*"
			}

			storageID := uuid.New()
			storage := &ent.StorageEntity{
				ID:            storageID,
				RefID:         order.PaymentID,
				FileName:      fileName,
				FilePath:      slipFilePath,
				FileSize:      slipFileSize,
				FileType:      fileType,
				IsActive:      true,
				RelatedEntity: ent.RelatedEntityPaymentFile,
				UploadedBy:    &uploadedBy,
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

		if err := s.upsertOrderPaymentReviewInTx(ctx, tx, order.ID, order.PaymentID, orderPaymentReviewStatusSubmitted, "", nil, nil); err != nil {
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

	review := new(ent.OrderPaymentReviewEntity)
	err := s.bunDB.DB().NewSelect().
		Model(review).
		Where("order_id = ?", orderID).
		OrderExpr("updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err == nil {
		switch strings.ToLower(strings.TrimSpace(review.ReviewStatus)) {
		case orderPaymentReviewStatusRejected:
			state.Rejected = true
			state.Reason = normalizePaymentRejectionReason(review.RejectedReason)
		case orderPaymentReviewStatusSubmitted, orderPaymentReviewStatusApproved:
			state.Submitted = true
		}
		return state, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return state, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
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

	isAppealApproval, err := s.isOrderPaymentAppealPendingReview(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	resolvedOrderStatus := ent.StatusTypePaid
	resolvedPaymentStatus := ent.PaymentTypeSuccess
	auditApprovedDetail := "Order payment approved by admin"
	cancellationReason := ""
	if isAppealApproval {
		resolvedOrderStatus = ent.StatusTypeCancelled
		resolvedPaymentStatus = ent.PaymentTypeRefunded
		auditApprovedDetail = "Order payment appeal approved by admin and refund approved"
		cancellationReason, err = s.getOrderPaymentAppealReason(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		if cancellationReason == "" {
			cancellationReason = "Approved payment appeal and refunded to customer"
		}
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		now := time.Now()

		payment := new(ent.PaymentEntity)
		if err := tx.NewSelect().Model(payment).Where("id = ?", order.PaymentID).Scan(ctx); err != nil {
			return err
		}

		payment.Status = resolvedPaymentStatus
		payment.ApprovedBy = &approverID
		payment.ApprovedAt = &now
		if _, err := tx.NewUpdate().Model(payment).Where("id = ?", payment.ID).Exec(ctx); err != nil {
			return err
		}

		previousStatus := order.Status
		order.Status = resolvedOrderStatus
		order.UpdatedAt = now

		if !isAppealApproval {
			if err := s.applyOrderStatusSideEffects(ctx, tx, order, previousStatus, approverID); err != nil {
				return err
			}
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
			ActionDetail: auditApprovedDetail,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		if isAppealApproval {
			if err := s.upsertOrderCancellationInTx(ctx, tx, order.ID, approverID, true, cancellationReason); err != nil {
				return err
			}

			statusLog := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.AuditActionUpdated,
				ActionType:   "order_status_transition",
				ActionID:     order.ID,
				ActionBy:     &approverID,
				Status:       ent.StatusAuditSuccesses,
				ActionDetail: fmt.Sprintf("Order status changed from %s to %s", previousStatus, order.Status),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(statusLog).Exec(ctx); err != nil {
				return err
			}
		}

		if err := s.upsertOrderPaymentReviewInTx(ctx, tx, order.ID, order.PaymentID, orderPaymentReviewStatusApproved, "", &approverID, &now); err != nil {
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
		OrderStatus:   string(resolvedOrderStatus),
		PaymentStatus: string(resolvedPaymentStatus),
		SlipAttached:  false,
	}, nil
}

func (s *Service) isOrderPaymentAppealPendingReview(ctx context.Context, orderID uuid.UUID) (bool, error) {
	latestLog := new(ent.AuditLogEntity)
	err := s.bunDB.DB().NewSelect().
		Model(latestLog).
		Where("action_id = ?", orderID).
		Where("action_type IN (?)", bun.In([]string{"order_payment_appealed", "order_payment_approved", "order_payment_rejected"})).
		Where("status = ?", ent.StatusAuditSuccesses).
		OrderExpr("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return latestLog.ActionType == "order_payment_appealed", nil
}

func (s *Service) RejectOrderPaymentService(ctx context.Context, orderID uuid.UUID, req *RejectOrderPaymentServiceRequest, approverID uuid.UUID) (*OrderPaymentServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.payment.reject.start`)

	reason := normalizePaymentRejectionReason(req.Reason)
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

		if err := s.upsertOrderPaymentReviewInTx(ctx, tx, order.ID, order.PaymentID, orderPaymentReviewStatusRejected, reason, &approverID, &now); err != nil {
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

func (s *Service) AppealOrderPaymentService(ctx context.Context, orderID uuid.UUID, req *AppealOrderPaymentServiceRequest, requesterID uuid.UUID) (*OrderPaymentServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.payment.appeal.start`)

	reason := normalizePaymentAppealReason(req.Reason)
	if reason == "" {
		return nil, errors.New("payment appeal reason is required")
	}

	order, err := s.ensureOrderAccess(ctx, orderID, requesterID, false)
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
	if !paymentReviewState.Rejected {
		return nil, errors.New("payment appeal is allowed only after rejection")
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		now := time.Now()

		payment := new(ent.PaymentEntity)
		if err := tx.NewSelect().Model(payment).Where("id = ?", order.PaymentID).Scan(ctx); err != nil {
			return err
		}

		payment.Status = ent.PaymentTypePending
		payment.ApprovedBy = nil
		payment.ApprovedAt = nil
		if _, err := tx.NewUpdate().Model(payment).Where("id = ?", payment.ID).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "order_payment_appealed",
			ActionID:     order.ID,
			ActionBy:     &requesterID,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Payment appeal reason: " + reason,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		if err := s.upsertOrderPaymentReviewInTx(ctx, tx, order.ID, order.PaymentID, orderPaymentReviewStatusSubmitted, "", nil, nil); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	span.AddEvent(`orders.svc.payment.appeal.success`)
	return &OrderPaymentServiceResponse{
		OrderID:       order.ID,
		PaymentID:     order.PaymentID,
		OrderStatus:   string(order.Status),
		PaymentStatus: string(ent.PaymentTypePending),
		SlipAttached:  false,
	}, nil
}

func parsePaymentRejectedReason(detail string) string {
	const prefix = "Payment rejected reason: "
	if strings.HasPrefix(detail, prefix) {
		return normalizePaymentRejectionReason(strings.TrimPrefix(detail, prefix))
	}
	return normalizePaymentRejectionReason(detail)
}

func parsePaymentAppealReason(detail string) string {
	const prefix = "Payment appeal reason: "
	if strings.HasPrefix(detail, prefix) {
		return normalizePaymentAppealReason(strings.TrimPrefix(detail, prefix))
	}
	return normalizePaymentAppealReason(detail)
}

func parseRefundRejectedReason(detail string) string {
	const prefix = "Refund rejected reason: "
	if strings.HasPrefix(detail, prefix) {
		return normalizeOrderCancellationReason(strings.TrimPrefix(detail, prefix))
	}
	return normalizeOrderCancellationReason(detail)
}

func normalizePaymentRejectionReason(reason string) string {
	normalized := strings.Join(strings.Fields(strings.TrimSpace(reason)), " ")
	return normalized
}

func normalizePaymentAppealReason(reason string) string {
	normalized := strings.Join(strings.Fields(strings.TrimSpace(reason)), " ")
	return normalized
}

func mapOrderStatusSummary(status ent.StatusTypeEnum, paymentSubmitted bool, paymentRejected bool) (string, string) {
	if status == ent.StatusTypePending {
		if paymentSubmitted {
			return "รอตรวจสอบการชำระเงิน", "รอแอดมินตรวจสอบและอนุมัติการชำระเงิน"
		}
		if paymentRejected {
			return "รอลูกค้าชำระเงินใหม่", "ชำระเงินใหม่และส่งหลักฐานการโอนอีกครั้ง"
		}
		return "รอชำระเงิน", "ชำระเงินและกดยืนยันการชำระเงิน"
	}

	switch status {
	case ent.StatusTypePaid:
		return "ชำระเงินแล้ว", "คำสั่งซื้อพร้อมจัดส่ง รอแอดมินเตรียมพัสดุ"
	case ent.StatusTypeRefundRequested:
		return "รอพิจารณาคืนเงิน", "รอแอดมินอนุมัติหรือปฏิเสธคำขอคืนเงิน"
	case ent.StatusTypeShipping:
		return "พร้อมติดตามพัสดุ", "ติดตามสถานะพัสดุจากเลขพัสดุที่ได้รับ"
	case ent.StatusTypeCompleted:
		return "คำสั่งซื้อสำเร็จ", "ได้รับสินค้าแล้ว สามารถสั่งซื้อครั้งถัดไปได้"
	case ent.StatusTypeCancelled:
		if paymentSubmitted {
			return "คืนเงินสำเร็จ", "แอดมินอนุมัติการคืนเงินเรียบร้อยแล้ว"
		}
		return "ยกเลิกคำสั่งซื้อ", "หากต้องการสินค้า กรุณาสร้างคำสั่งซื้อใหม่"
	default:
		return string(status), ""
	}
}

func (s *Service) getOrderShippingTrackingNo(ctx context.Context, orderID uuid.UUID) (string, error) {
	tracking := new(ent.OrderShippingTrackingEntity)
	err := s.bunDB.DB().NewSelect().
		Model(tracking).
		Where("order_id = ?", orderID).
		OrderExpr("updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err == nil {
		return strings.TrimSpace(tracking.TrackingNo), nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	return "", nil
}

func (s *Service) upsertOrderShippingTrackingInTx(ctx context.Context, tx bun.Tx, orderID uuid.UUID, trackingNo string, updatedBy *uuid.UUID) error {
	normalizedTrackingNo := strings.TrimSpace(trackingNo)
	if normalizedTrackingNo == "" {
		return nil
	}

	now := time.Now()
	record := &ent.OrderShippingTrackingEntity{
		ID:         uuid.New(),
		OrderID:    orderID,
		TrackingNo: normalizedTrackingNo,
		UpdatedBy:  updatedBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if _, err := tx.NewInsert().
		Model(record).
		On("CONFLICT (order_id) DO UPDATE").
		Set("tracking_no = EXCLUDED.tracking_no").
		Set("updated_by = EXCLUDED.updated_by").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) upsertOrderPaymentReviewInTx(
	ctx context.Context,
	tx bun.Tx,
	orderID uuid.UUID,
	paymentID uuid.UUID,
	reviewStatus string,
	rejectedReason string,
	reviewedBy *uuid.UUID,
	reviewedAt *time.Time,
) error {
	normalizedStatus := strings.ToLower(strings.TrimSpace(reviewStatus))
	if normalizedStatus == "" {
		normalizedStatus = orderPaymentReviewStatusSubmitted
	}

	normalizedReason := normalizePaymentRejectionReason(rejectedReason)
	now := time.Now()

	record := &ent.OrderPaymentReviewEntity{
		ID:             uuid.New(),
		OrderID:        orderID,
		PaymentID:      paymentID,
		ReviewStatus:   normalizedStatus,
		RejectedReason: normalizedReason,
		ReviewedBy:     reviewedBy,
		ReviewedAt:     reviewedAt,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if _, err := tx.NewInsert().
		Model(record).
		On("CONFLICT (order_id) DO UPDATE").
		Set("payment_id = EXCLUDED.payment_id").
		Set("review_status = EXCLUDED.review_status").
		Set("rejected_reason = EXCLUDED.rejected_reason").
		Set("reviewed_by = EXCLUDED.reviewed_by").
		Set("reviewed_at = EXCLUDED.reviewed_at").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) getOrderCancellationReason(ctx context.Context, orderID uuid.UUID) (string, error) {
	record := new(ent.OrderCancellationEntity)
	err := s.bunDB.DB().NewSelect().
		Model(record).
		Where("order_id = ?", orderID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return normalizeOrderCancellationReason(record.Reason), nil
}

func (s *Service) getOrderRefundRejectionReason(ctx context.Context, orderID uuid.UUID) (string, error) {
	latestLog := new(ent.AuditLogEntity)
	err := s.bunDB.DB().NewSelect().
		Model(latestLog).
		Where("action_id = ?", orderID).
		Where("action_type = ?", "order_refund_rejected").
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

	return parseRefundRejectedReason(latestLog.ActionDetail), nil
}

func (s *Service) getOrderPaymentAppealReason(ctx context.Context, orderID uuid.UUID) (string, error) {
	latestLog := new(ent.AuditLogEntity)
	err := s.bunDB.DB().NewSelect().
		Model(latestLog).
		Where("action_id = ?", orderID).
		Where("action_type = ?", "order_payment_appealed").
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

	return parsePaymentAppealReason(latestLog.ActionDetail), nil
}

func parseShippingTrackingNumber(detail string) string {
	const prefix = "Shipping tracking number: "
	if strings.HasPrefix(detail, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(detail, prefix))
	}
	return strings.TrimSpace(detail)
}

func mapNotificationTitleMessage(actionType string, actionDetail string, orderNo string) (string, string) {
	orderRef := strings.TrimSpace(orderNo)
	if orderRef == "" {
		orderRef = "คำสั่งซื้อของคุณ"
	} else {
		orderRef = "คำสั่งซื้อ " + orderRef
	}

	switch actionType {
	case "order_payment_appealed":
		reason := parsePaymentAppealReason(actionDetail)
		if reason == "" {
			reason = "ลูกค้าส่งคำขออุทธรณ์การชำระเงิน"
		}
		return "อุทธรณ์การชำระเงิน", orderRef + " ส่งคำขออุทธรณ์การชำระเงิน: " + reason
	case "order_payment_approved":
		return "ชำระเงินได้รับอนุมัติ", orderRef + " ได้รับการอนุมัติการชำระเงินแล้ว และกำลังเข้าสู่ขั้นตอนจัดส่ง"
	case "order_payment_rejected":
		reason := parsePaymentRejectedReason(actionDetail)
		if reason == "" {
			reason = "กรุณาชำระเงินใหม่และส่งหลักฐานอีกครั้ง"
		}
		return "ไม่อนุมัติการชำระเงิน", orderRef + " ไม่ผ่านการอนุมัติ: " + reason
	case "order_refund_rejected":
		reason := parseRefundRejectedReason(actionDetail)
		if reason == "" {
			reason = "แอดมินปฏิเสธคำขอคืนเงิน"
		}
		return "ไม่อนุมัติการคืนเงิน", orderRef + " ถูกปฏิเสธการคืนเงิน: " + reason
	case "order_shipping_tracking_updated":
		trackingNo := parseShippingTrackingNumber(actionDetail)
		if trackingNo == "" {
			return "อัปเดตการจัดส่ง", orderRef + " มีการอัปเดตข้อมูลการจัดส่ง"
		}
		return "อัปเดตเลขพัสดุ", orderRef + " อัปเดตเลขพัสดุ: " + trackingNo
	case "order_status_transition":
		fromStatus, toStatus := parseOrderStatusTransitionDetail(actionDetail)
		if toStatus == string(ent.StatusTypeRefundRequested) {
			return "ส่งคำขอคืนเงิน", orderRef + " ส่งคำขอคืนเงินและอยู่ระหว่างรอพิจารณา"
		}
		if fromStatus == string(ent.StatusTypeRefundRequested) && toStatus == string(ent.StatusTypeCancelled) {
			return "อนุมัติคืนเงินแล้ว", orderRef + " ได้รับอนุมัติคืนเงินเรียบร้อยแล้ว"
		}
		if toStatus == string(ent.StatusTypeShipping) {
			return "คำสั่งซื้อกำลังจัดส่ง", orderRef + " เปลี่ยนสถานะเป็นกำลังจัดส่งแล้ว"
		}
		if toStatus == string(ent.StatusTypeCompleted) {
			return "คำสั่งซื้อสำเร็จ", orderRef + " จัดส่งสำเร็จเรียบร้อยแล้ว"
		}
		if toStatus == string(ent.StatusTypeCancelled) {
			return "คำสั่งซื้อถูกยกเลิก", orderRef + " ถูกยกเลิกแล้ว"
		}
		if fromStatus != "" && toStatus != "" {
			return "สถานะคำสั่งซื้ออัปเดต", orderRef + " เปลี่ยนสถานะจาก " + fromStatus + " เป็น " + toStatus
		}
		return "สถานะคำสั่งซื้ออัปเดต", orderRef + " มีการอัปเดตสถานะ"
	default:
		return "อัปเดตคำสั่งซื้อ", orderRef + " มีการอัปเดตใหม่"
	}
}
