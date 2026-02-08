package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListPaymentsRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
