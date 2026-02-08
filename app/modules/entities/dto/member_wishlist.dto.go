package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberWishlistRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
