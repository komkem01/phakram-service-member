package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListOrderItemsRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
