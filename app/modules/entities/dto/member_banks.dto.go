package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberBanksRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
