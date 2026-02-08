package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMembersRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
