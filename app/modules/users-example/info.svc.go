package users

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"log/slog"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/config/i18n"
// 	"time"

// 	"github.com/google/uuid"
// )

// type InfoUserService struct {
// 	ID        uuid.UUID
// 	Username  string
// 	FirstName string
// 	LastName  string
// 	Email     string
// 	Phone     string
// 	Status    string
// 	Role      string
// 	CreatedAt string
// 	UpdatedAt string
// }

// type InfoServiceResponse struct {
// 	ID        uuid.UUID `json:"id"`
// 	Username  string    `json:"username"`
// 	FirstName string    `json:"first_name"`
// 	LastName  string    `json:"last_name"`
// 	Email     string    `json:"email"`
// 	Phone     string    `json:"phone"`
// 	Status    string    `json:"status"`
// 	Role      string    `json:"role"`
// 	CreatedAt string    `json:"created_at"`
// 	UpdatedAt string    `json:"updated_at"`
// }

// func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoServiceResponse, error) {
// 	span, log := utils.LogSpanFromContext(ctx)
// 	span.AddEvent(`user.svc.info.start`)

// 	data, err := s.db.GetUserByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, i18n.ErrUserNotFound
// 		}
// 		log.With(slog.Any(`body`, id)).Errf(`internal: %s`, err)
// 		return nil, i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`user.svc.info.db`)
// 	email, err := s.dbContactEmail.GetContactEmailByUserID(ctx, data.ID)
// 	if err != nil {
// 		if !errors.Is(err, sql.ErrNoRows) {
// 			log.With(slog.Any(`body`, id)).Errf(`internal: %s`, err)
// 			return nil, i18n.ErrEmailNotFound
// 		}
// 		return nil, i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`user.svc.info.db_email`)
// 	phone, err := s.dbContactPhone.GetContactPhoneByUserID(ctx, data.ID)
// 	if err != nil {
// 		if !errors.Is(err, sql.ErrNoRows) {
// 			log.With(slog.Any(`body`, id)).Errf(`internal: %s`, err)
// 			return nil, i18n.ErrPhoneNotFound
// 		}
// 		return nil, i18n.ErrInternalServerError
// 	}

// 	response := &InfoServiceResponse{
// 		ID:        data.ID,
// 		Username:  data.Username,
// 		FirstName: data.FirstName,
// 		LastName:  data.LastName,
// 		Email:     email.Email,
// 		Phone:     phone.PhoneNumber,
// 		Status:    string(data.Status),
// 		Role:      string(data.Role),
// 		CreatedAt: data.CreatedAt.Format(time.RFC3339),
// 		UpdatedAt: data.UpdatedAt.Format(time.RFC3339),
// 	}
// 	span.AddEvent(`user.svc.info.copy`)
// 	return response, nil
// }
