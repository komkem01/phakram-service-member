package member_files

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type ListMemberFileServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListMemberFileServiceResponses struct {
	ID        uuid.UUID `json:"id"`
	MemberID  uuid.UUID `json:"member_id"`
	FileID    uuid.UUID `json:"file_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListMemberFileServiceRequest) ([]*ListMemberFileServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_files.svc.list.start`)

	data, page, err := s.db.ListMemberFiles(ctx, &entitiesdto.ListMemberFilesRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListMemberFileServiceResponses
	for _, item := range data {
		temp := &ListMemberFileServiceResponses{
			ID:        item.ID,
			MemberID:  item.MemberID,
			FileID:    item.FileID,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`member_files.svc.list.copy`)
	return response, page, nil
}
