package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListCartItemsRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
