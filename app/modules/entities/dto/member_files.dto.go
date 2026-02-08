package entitiesdto

import (
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberFilesRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}
