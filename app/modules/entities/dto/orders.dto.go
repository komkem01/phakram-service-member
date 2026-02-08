package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListOrdersRequest struct {
	base.RequestPaginate
	MemberID  uuid.UUID
	Search    string
	Status    string
	StartDate int64
	EndDate   int64
}
