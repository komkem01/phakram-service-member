package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrderShippingTrackingEntity struct {
	bun.BaseModel `bun:"table:order_shipping_trackings"`

	ID         uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	OrderID    uuid.UUID  `bun:"order_id,type:uuid" json:"order_id"`
	TrackingNo string     `bun:"tracking_no" json:"tracking_no"`
	UpdatedBy  *uuid.UUID `bun:"updated_by,type:uuid" json:"updated_by"`
	CreatedAt  time.Time  `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
