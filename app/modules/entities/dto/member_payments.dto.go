package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberPaymentsRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
