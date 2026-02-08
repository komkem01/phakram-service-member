package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateOrderService struct {
	CartID         uuid.UUID   `json:"cart_id"`
	CartItemIDs    []uuid.UUID `json:"cart_item_ids"`
	OrderNo        string      `json:"order_no"`
	MemberID       uuid.UUID   `json:"member_id"`
	PaymentID      uuid.UUID   `json:"payment_id"`
	AddressID      uuid.UUID   `json:"address_id"`
	Address        *CreateOrderAddress
	DiscountAmount decimal.Decimal `json:"discount_amount"`
}

type CreateOrderAddress struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	IsDefault     bool   `json:"is_default"`
	AddressNo     string `json:"address_no"`
	Village       string `json:"village"`
	Alley         string `json:"alley"`
	SubDistrictID string `json:"sub_district_id"`
	DistrictID    string `json:"district_id"`
	ProvinceID    string `json:"province_id"`
	ZipcodeID     string `json:"zipcode_id"`
}

func (s *Service) CreateOrderService(ctx context.Context, req *CreateOrderService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.create.start`)

	orderNo, err := utils.GenerateOrderNo()
	if err != nil {
		return err
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		cart := new(ent.CartEntity)
		if err := tx.NewSelect().
			Model(cart).
			Where("id = ?", req.CartID).
			Limit(1).
			Scan(ctx); err != nil {
			if err == sql.ErrNoRows {
				return errors.New("cart not found")
			}
			return err
		}
		if cart.MemberID != req.MemberID {
			return errors.New("cart not found")
		}

		addressID := req.AddressID
		if addressID != uuid.Nil {
			address := new(ent.MemberAddressEntity)
			if err := tx.NewSelect().
				Model(address).
				Where("id = ?", addressID).
				Where("member_id = ?", req.MemberID).
				Limit(1).
				Scan(ctx); err != nil {
				if err == sql.ErrNoRows {
					return errors.New("address not found")
				}
				return err
			}
		} else if req.Address != nil {
			subDistrictID, err := uuid.Parse(req.Address.SubDistrictID)
			if err != nil {
				return err
			}
			districtID, err := uuid.Parse(req.Address.DistrictID)
			if err != nil {
				return err
			}
			provinceID, err := uuid.Parse(req.Address.ProvinceID)
			if err != nil {
				return err
			}
			zipcodeID, err := uuid.Parse(req.Address.ZipcodeID)
			if err != nil {
				return err
			}

			address := &ent.MemberAddressEntity{
				ID:            uuid.New(),
				MemberID:      req.MemberID,
				FirstName:     req.Address.FirstName,
				LastName:      req.Address.LastName,
				Phone:         req.Address.Phone,
				IsDefault:     req.Address.IsDefault,
				AddressNo:     req.Address.AddressNo,
				Village:       req.Address.Village,
				Alley:         req.Address.Alley,
				SubDistrictID: subDistrictID,
				DistrictID:    districtID,
				ProvinceID:    provinceID,
				ZipcodeID:     zipcodeID,
			}
			if _, err := tx.NewInsert().Model(address).Exec(ctx); err != nil {
				return err
			}
			addressID = address.ID

			actionBy := req.MemberID
			now := time.Now()
			auditLog := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.ActionAuditCreate,
				ActionType:   "member_address",
				ActionID:     &address.ID,
				ActionBy:     &actionBy,
				Status:       ent.StatusAuditSuccess,
				ActionDetail: "Created member address " + address.ID.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("missing address")
		}

		items := make([]*ent.CartItemEntity, 0)
		selectItems := tx.NewSelect().
			Model(&items).
			Where("cart_id = ?", req.CartID)
		if len(req.CartItemIDs) > 0 {
			selectItems = selectItems.Where("id IN (?)", bun.In(req.CartItemIDs))
		}
		if err := selectItems.Scan(ctx); err != nil {
			return err
		}
		if len(items) == 0 {
			return errors.New("cart items not found")
		}

		totalAmount := decimal.Zero
		now := time.Now()
		for _, item := range items {
			itemTotal := item.PricePerUnit.Mul(decimal.NewFromInt(int64(item.Quantity)))
			totalAmount = totalAmount.Add(itemTotal)

			stock := new(ent.ProductStockEntity)
			if err := tx.NewSelect().
				Model(stock).
				Where("product_id = ?", item.ProductID).
				Limit(1).
				Scan(ctx); err != nil {
				if err == sql.ErrNoRows {
					return errors.New("product stock not found")
				}
				return err
			}
			if stock.Remaining < item.Quantity {
				return fmt.Errorf("product stock is not enough")
			}
			stock.Remaining = stock.Remaining - item.Quantity
			if _, err := tx.NewUpdate().
				Model(stock).
				Set("remaining = ?", stock.Remaining).
				Where("id = ?", stock.ID).
				Exec(ctx); err != nil {
				return err
			}

			actionBy := req.MemberID
			stockAudit := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.ActionAuditUpdate,
				ActionType:   "product_stock",
				ActionID:     &stock.ID,
				ActionBy:     &actionBy,
				Status:       ent.StatusAuditSuccess,
				ActionDetail: "Updated product stock " + stock.ID.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(stockAudit).Exec(ctx); err != nil {
				return err
			}
		}

		netAmount := totalAmount.Sub(req.DiscountAmount)
		order := &ent.OrderEntity{
			ID:             uuid.New(),
			OrderNo:        orderNo,
			MemberID:       req.MemberID,
			PaymentID:      req.PaymentID,
			AddressID:      addressID,
			Status:         ent.StatusTypePending,
			TotalAmount:    totalAmount,
			DiscountAmount: req.DiscountAmount,
			NetAmount:      netAmount,
		}
		insertOrder := tx.NewInsert().Model(order)
		if req.PaymentID == uuid.Nil {
			insertOrder = insertOrder.ExcludeColumn("payment_id")
		}
		if _, err := insertOrder.Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		orderAudit := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "order",
			ActionID:     &order.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created order " + order.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(orderAudit).Exec(ctx); err != nil {
			return err
		}

		itemIDs := make([]uuid.UUID, 0, len(items))
		for _, item := range items {
			orderItem := &ent.OrderItemEntity{
				ID:              uuid.New(),
				OrderID:         order.ID,
				ProductID:       item.ProductID,
				Quantity:        item.Quantity,
				PricePerUnit:    item.PricePerUnit,
				TotalItemAmount: item.PricePerUnit.Mul(decimal.NewFromInt(int64(item.Quantity))),
			}
			if _, err := tx.NewInsert().Model(orderItem).Exec(ctx); err != nil {
				return err
			}

			itemAudit := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.ActionAuditCreate,
				ActionType:   "order_item",
				ActionID:     &orderItem.ID,
				ActionBy:     &actionBy,
				Status:       ent.StatusAuditSuccess,
				ActionDetail: "Created order item " + orderItem.ID.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(itemAudit).Exec(ctx); err != nil {
				return err
			}

			itemIDs = append(itemIDs, item.ID)
		}

		if _, err := tx.NewDelete().
			Model((*ent.CartItemEntity)(nil)).
			Where("id IN (?)", bun.In(itemIDs)).
			Exec(ctx); err != nil {
			return err
		}

		remainingCount, err := tx.NewSelect().
			Model((*ent.CartItemEntity)(nil)).
			Where("cart_id = ?", req.CartID).
			Count(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if remainingCount == 0 {
			if _, err := tx.NewUpdate().
				Model((*ent.CartEntity)(nil)).
				Set("is_active = ?", false).
				Where("id = ?", req.CartID).
				Exec(ctx); err != nil {
				return err
			}

			cartAudit := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.ActionAuditUpdate,
				ActionType:   "cart",
				ActionID:     &req.CartID,
				ActionBy:     &actionBy,
				Status:       ent.StatusAuditSuccess,
				ActionDetail: "Deactivated cart " + req.CartID.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(cartAudit).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}
	span.AddEvent(`orders.svc.create.success`)
	return nil
}
