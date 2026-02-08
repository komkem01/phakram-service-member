package users

// import (
// 	"context"
// 	"phakram-craft/app/modules/entities/ent"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/hashing"
// 	"phakram-craft/config/i18n"

// 	"github.com/google/uuid"
// )

// type ChangePasswordService struct {
// 	OldPassword string `json:"old_password"`
// 	NewPassword string `json:"new_password"`
// }

// func (s *Service) ChangePasswordService(ctx context.Context, userID uuid.UUID, req *ChangePasswordService) error {
// 	span, log := utils.LogSpanFromContext(ctx)
// 	span.AddEvent(`users.svc.change_password.start`)

// 	// Get user to verify existence
// 	_, err := s.db.GetUserByID(ctx, userID)
// 	if err != nil {
// 		span.AddEvent(`users.svc.change_password.user_not_found`)
// 		log.Errf("User not found: %v", err)
// 		return i18n.ErrUserNotFound
// 	}
// 	span.AddEvent(`users.svc.change_password.user_found`)

// 	// Get user credentials (email auth)
// 	credentials, err := s.dbCredential.GetCredentialsByUserID(ctx, userID)
// 	if err != nil || len(credentials) == 0 {
// 		span.AddEvent(`users.svc.change_password.credentials_not_found`)
// 		log.Errf("Credentials not found for user: %v", err)
// 		return i18n.ErrUserNotFound
// 	}

// 	// Find email credential with password
// 	var emailCred *ent.UserCredentialEntity
// 	for _, cred := range credentials {
// 		if cred.AuthType == ent.AuthTypeEmail && cred.PasswordHash != nil {
// 			emailCred = cred
// 			break
// 		}
// 	}

// 	if emailCred == nil {
// 		span.AddEvent(`users.svc.change_password.email_credential_not_found`)
// 		log.Warnf("Email credential not found for user: %s", userID)
// 		return i18n.ErrInvalidPassword
// 	}
// 	span.AddEvent(`users.svc.change_password.credential_found`)

// 	// Verify old password
// 	if !hashing.CheckPasswordHash([]byte(*emailCred.PasswordHash), []byte(req.OldPassword)) {
// 		span.AddEvent(`users.svc.change_password.invalid_old_password`)
// 		log.Warnf("Invalid old password for user: %s", userID)
// 		return i18n.ErrInvalidPassword
// 	}
// 	span.AddEvent(`users.svc.change_password.old_password_verified`)

// 	// Hash new password
// 	hashedPassword, err := hashing.HashPassword(req.NewPassword)
// 	if err != nil {
// 		span.AddEvent(`users.svc.change_password.error_hashing_password`)
// 		log.Errf("Error hashing password: %v", err)
// 		return i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`users.svc.change_password.new_password_hashed`)

// 	// Update password in credentials
// 	if err := s.dbCredential.UpdatePasswordHash(ctx, emailCred.ID, string(hashedPassword)); err != nil {
// 		span.AddEvent(`users.svc.change_password.error_updating_credential`)
// 		log.Errf("Error updating credential password: %v", err)
// 		return i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`users.svc.change_password.success`)

// 	return nil
// }
