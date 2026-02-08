package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberAddressesRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
