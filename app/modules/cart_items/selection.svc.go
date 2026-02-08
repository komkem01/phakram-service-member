package cart_items

import (
	"context"
	"database/sql"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SelectCartItemsServiceRequest struct {
	CartID   uuid.UUID
	ItemIDs  []uuid.UUID
	MemberID uuid.UUID
	IsAdmin  bool
}

func (s *Service) SelectItemsService(ctx context.Context, req *SelectCartItemsServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`cart_items.svc.select.start`)

	if !req.IsAdmin {
		var exists int
		err := s.bunDB.DB().NewSelect().
			Table("carts").
			ColumnExpr("1").
			Where("id = ? AND member_id = ?", req.CartID, req.MemberID).
			Limit(1).
			Scan(ctx, &exists)
		if err != nil {
			if err == sql.ErrNoRows {
				return sql.ErrNoRows
			}
			return err
		}
	}

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().
			Table("cart_items").
			Set("is_selected = ?", false).
			Where("cart_id = ?", req.CartID).
			Exec(ctx); err != nil {
			return err
		}

		if len(req.ItemIDs) == 0 {
			return nil
		}

		if _, err := tx.NewUpdate().
			Table("cart_items").
			Set("is_selected = ?", true).
			Where("cart_id = ?", req.CartID).
			Where("id IN (?)", bun.In(req.ItemIDs)).
			Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	span.AddEvent(`cart_items.svc.select.success`)
	return nil
}
