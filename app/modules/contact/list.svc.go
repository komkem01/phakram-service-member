package contact

import (
	"context"
	"strings"

	"phakram/app/utils/base"

	"github.com/uptrace/bun"
)

type ListContactServiceRequest struct {
	base.RequestPaginate
	SendStatus string
	ReadStatus string
}

type ListContactServiceResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Subject    string  `json:"subject"`
	Message    string  `json:"message"`
	SendStatus string  `json:"send_status"`
	IsRead     bool    `json:"is_read"`
	SendError  *string `json:"send_error"`
	SentAt     *string `json:"sent_at"`
	ReadAt     *string `json:"read_at"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func (s *Service) List(ctx context.Context, req *ListContactServiceRequest) ([]*ListContactServiceResponse, *base.ResponsePaginate, error) {
	query := s.bunDB.DB().NewSelect().Model((*contactMessageRecord)(nil))

	status := strings.TrimSpace(strings.ToLower(req.SendStatus))
	if status != "" {
		query = query.Where("LOWER(send_status) = ?", status)
	}

	readStatus := strings.TrimSpace(strings.ToLower(req.ReadStatus))
	if readStatus == "read" {
		query = query.Where("is_read = TRUE")
	}
	if readStatus == "unread" {
		query = query.Where("is_read = FALSE")
	}

	search := strings.TrimSpace(req.Search)
	if search != "" {
		query = query.Where("(name ILIKE ? OR email ILIKE ? OR subject ILIKE ?)", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	sortBy := strings.ToLower(strings.TrimSpace(req.SortBy))
	switch sortBy {
	case "created_at", "updated_at", "sent_at", "name", "email", "send_status":
	default:
		sortBy = "created_at"
	}
	if req.OrderBy == "" {
		req.OrderBy = "desc"
	}

	orderBy := "ASC"
	if strings.EqualFold(req.OrderBy, "desc") {
		orderBy = "DESC"
	}

	items := make([]*contactMessageRecord, 0)
	query = s.bunDB.DB().NewSelect().
		Model(&items).
		OrderExpr("? "+orderBy, bun.Ident(sortBy))

	if status != "" {
		query = query.Where("LOWER(send_status) = ?", status)
	}
	if readStatus == "read" {
		query = query.Where("is_read = TRUE")
	}
	if readStatus == "unread" {
		query = query.Where("is_read = FALSE")
	}
	if search != "" {
		query = query.Where("(name ILIKE ? OR email ILIKE ? OR subject ILIKE ?)", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	req.SetOffsetLimit(query)
	if err := query.Scan(ctx); err != nil {
		return nil, nil, err
	}

	resp := make([]*ListContactServiceResponse, 0, len(items))
	for _, item := range items {
		var sentAt *string
		if item.SentAt != nil {
			v := item.SentAt.Format("2006-01-02T15:04:05Z07:00")
			sentAt = &v
		}

		var readAt *string
		if item.ReadAt != nil {
			v := item.ReadAt.Format("2006-01-02T15:04:05Z07:00")
			readAt = &v
		}

		var sendError *string
		if strings.TrimSpace(item.SendError) != "" {
			v := item.SendError
			sendError = &v
		}

		resp = append(resp, &ListContactServiceResponse{
			ID:         item.ID.String(),
			Name:       item.Name,
			Email:      item.Email,
			Subject:    item.Subject,
			Message:    item.Message,
			SendStatus: item.SendStatus,
			IsRead:     item.IsRead,
			SendError:  sendError,
			SentAt:     sentAt,
			ReadAt:     readAt,
			CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return resp, &base.ResponsePaginate{
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: int64(total),
	}, nil
}
