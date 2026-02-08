package users

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"log/slog"
// 	"phakram-craft/app/modules/entities/ent"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/config/i18n"

// 	"github.com/google/uuid"
// )

// func (s *Service) DeleteService(ctx context.Context, actorRole ent.UserRole, id uuid.UUID) error {
// 	span, log := utils.LogSpanFromContext(ctx)
// 	span.AddEvent(`users.svc.delete.start`)

// 	user, err := s.db.GetUserByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return i18n.ErrUserNotFound
// 		}
// 		log.With(slog.Any(`body`, id)).Errf(`internal: %s`, err)
// 		return i18n.ErrInternalServerError
// 	}

// 	switch actorRole {
// 	case ent.UserRoleSuperadmin:
// 		if user.Role == ent.UserRoleSuperadmin {
// 			return i18n.ErrForbidden
// 		}
// 	case ent.UserRoleAdmin:
// 		if user.Role == ent.UserRoleSuperadmin {
// 			return i18n.ErrForbidden
// 		}
// 	default:
// 		return i18n.ErrForbidden
// 	}

// 	if err := s.db.DeleteUser(ctx, id); err != nil {
// 		log.With(slog.Any(`body`, id)).Errf(`internal: %s`, err)
// 		return i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`users.svc.delete.db`)
// 	return nil
// }
