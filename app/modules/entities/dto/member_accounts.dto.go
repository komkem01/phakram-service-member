package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberAccountsRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
