package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberAddressEntity struct {
	bun.BaseModel `bun:"table:member_addresses"`

	ID            uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	MemberID      uuid.UUID  `bun:"member_id,type:uuid" json:"member_id"`
	FirstName     string     `bun:"first_name" json:"first_name"`
	LastName      string     `bun:"last_name" json:"last_name"`
	Phone         string     `bun:"phone" json:"phone"`
	IsDefault     bool       `bun:"is_default" json:"is_default"`
	AddressNo     string     `bun:"address_no" json:"address_no"`
	Village       string     `bun:"village" json:"village"`
	Alley         string     `bun:"alley" json:"alley"`
	SubDistrictID uuid.UUID  `bun:"sub_district_id,type:uuid" json:"sub_district_id"`
	DistrictID    uuid.UUID  `bun:"district_id,type:uuid" json:"district_id"`
	ProvinceID    uuid.UUID  `bun:"province_id,type:uuid" json:"province_id"`
	ZipcodeID     uuid.UUID  `bun:"zipcode_id,type:uuid" json:"zipcode_id"`
	CreatedAt     time.Time  `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at"`
}
