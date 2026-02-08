package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListCartsRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
