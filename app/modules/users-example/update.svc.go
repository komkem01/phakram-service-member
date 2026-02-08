package users

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"log/slog"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/config/i18n"
// 	"phakram-craft/internal/log"

// 	"github.com/google/uuid"
// )

// type UpdateUserServiceRequestUri struct {
// 	ID uuid.UUID
// }

// type UpdateUserService struct {
// 	Username  string `json:"username"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// }

// func (s *Service) UpdateUserService(ctx context.Context, id uuid.UUID, req *UpdateUserService) error {
// 	span, _ := utils.LogSpanFromContext(ctx)
// 	span.AddEvent(`users.svc.update.start`)

// 	user, err := s.db.GetUserByID(ctx, id)
// 	if err != nil {
// 		if !errors.Is(err, sql.ErrNoRows) {
// 			log.With(slog.Any(`body`, id)).Errf(`internal: %s`, err)
// 			return i18n.ErrUserNotFound
// 		}
// 		return i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`users.svc.update.user_found`)

// 	user.Username = req.Username
// 	user.FirstName = req.FirstName
// 	user.LastName = req.LastName

// 	if err := s.db.UpdateUser(ctx, user); err != nil {
// 		span.AddEvent(`users.svc.update.error_updating_user`)
// 		return i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`users.svc.update.user_updated`)

// 	return nil
// }
