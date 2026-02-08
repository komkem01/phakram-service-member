package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberTransactionsRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
