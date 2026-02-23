package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type promotionUsageRecord struct {
	bun.BaseModel `bun:"table:promotion_usages"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid"`
	PromotionID    uuid.UUID  `bun:"promotion_id,type:uuid,notnull"`
	MemberID       uuid.UUID  `bun:"member_id,type:uuid,notnull"`
	OrderID        *uuid.UUID `bun:"order_id,type:uuid"`
	DiscountAmount float64    `bun:"discount_amount,notnull"`
	UsedAt         time.Time  `bun:"used_at,notnull"`
	CreatedAt      time.Time  `bun:"created_at,notnull"`
}
