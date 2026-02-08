package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListPaymentFilesRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
