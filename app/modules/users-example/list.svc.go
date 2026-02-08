package users

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"log/slog"
// 	entitiesdto "phakram-craft/app/modules/entities/dto"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/base"

// 	"github.com/google/uuid"
// )

// type ListServiceRequest struct {
// 	base.RequestPaginate
// 	Role string
// }

// type ListServiceResponse struct {
// 	ID        uuid.UUID `json:"id"`
// 	Username  string    `json:"username"`
// 	FirstName string    `json:"first_name"`
// 	LastName  string    `json:"last_name"`
// 	Status    string    `json:"status"`
// 	Role      string    `json:"role"`
// 	Email     string    `json:"email"`
// 	Phone     string    `json:"phone"`
// 	CreatedAt string    `json:"created_at"`
// 	UpdatedAt string    `json:"updated_at"`
// }

// func (s *Service) ListService(ctx context.Context, req *ListServiceRequest) ([]*ListServiceResponse, *base.ResponsePaginate, error) {
// 	span, log := utils.LogSpanFromContext(ctx)
// 	span.AddEvent(`users.svc.list.start`)

// 	data, page, err := s.db.ListUsers(ctx, &entitiesdto.ListUsersRequest{
// 		RequestPaginate: req.RequestPaginate,
// 		Role:            req.Role,
// 	})
// 	if err != nil {
// 		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 		return nil, nil, err
// 	}
// 	var response []*ListServiceResponse
// 	for _, item := range data {
// 		email, err := s.dbContactEmail.GetContactEmailByUserID(ctx, item.ID)
// 		if err != nil {
// 			if !errors.Is(err, sql.ErrNoRows) {
// 				log.With(slog.Any(`body`, item.ID)).Errf(`internal: %s`, err)
// 			} else {
// 				log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 				return nil, nil, err
// 			}
// 		}
// 		phone, err := s.dbContactPhone.GetContactPhoneByUserID(ctx, item.ID)
// 		if err != nil {
// 			if !errors.Is(err, sql.ErrNoRows) {
// 				log.With(slog.Any(`body`, item.ID)).Errf(`internal: %s`, err)
// 			} else {
// 				log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
// 				return nil, nil, err
// 			}
// 		}

// 		temp := &ListServiceResponse{
// 			ID:        item.ID,
// 			Username:  item.Username,
// 			FirstName: item.FirstName,
// 			LastName:  item.LastName,
// 			Status:    string(item.Status),
// 			Role:      string(item.Role),
// 			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 		}
// 		if email != nil {
// 			temp.Email = email.Email
// 		}
// 		if phone != nil {
// 			temp.Phone = phone.PhoneNumber
// 		}
// 		response = append(response, temp)
// 	}
// 	span.AddEvent(`users.svc.list.copy`)
// 	return response, page, nil
// }
