package users

// import (
// 	"context"
// 	"phakram-craft/app/modules/entities/ent"
// 	"phakram-craft/app/utils"
// 	"phakram-craft/app/utils/hashing"
// 	"phakram-craft/config/i18n"

// 	"github.com/google/uuid"
// 	"github.com/uptrace/bun"
// )

// type CreateUserService struct {
// 	ID        uuid.UUID `json:"id"`
// 	Username  string    `json:"username"`
// 	Password  string    `json:"password"`
// 	FirstName string    `json:"first_name"`
// 	LastName  string    `json:"last_name"`
// 	Email     string    `json:"email"`
// 	Phone     string    `json:"phone"`
// 	Role      ent.UserRole
// }

// func (s *Service) CreateUserService(ctx context.Context, req *CreateUserService) error {
// 	span, _ := utils.LogSpanFromContext(ctx)
// 	span.AddEvent(`users.svc.create.start`)

// 	_, err := s.db.GetUserByUsername(ctx, req.Username)
// 	if err == nil {
// 		span.AddEvent(`users.svc.create.username_exists`)
// 		return i18n.ErrUsernameExists
// 	}

// 	// Hash password
// 	hashedPassword, err := hashing.HashPassword(req.Password)
// 	if err != nil {
// 		span.AddEvent(`users.svc.create.error_hashing_password`)
// 		return i18n.ErrInternalServerError
// 	}
// 	span.AddEvent(`users.svc.create.password_hashed`)

// 	id := uuid.New()

// 	role := req.Role
// 	if role == "" {
// 		role = ent.UserRoleUser
// 	}

// 	// Use transaction to ensure atomic operation
// 	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
// 		// Create user
// 		user := &ent.UserEntity{
// 			ID:        id,
// 			Username:  req.Username,
// 			FirstName: req.FirstName,
// 			LastName:  req.LastName,
// 			Status:    ent.UserStatusActive,
// 			Role:      role,
// 		}
// 		if err := s.db.CreateUser(ctx, user); err != nil {
// 			span.AddEvent(`users.svc.create.error_creating_user`)
// 			return err
// 		}
// 		span.AddEvent(`users.svc.create.user_created`)

// 		// Create user credential (email + password)
// 		passwordHashStr := string(hashedPassword)
// 		credential := &ent.UserCredentialEntity{
// 			ID:           uuid.New(),
// 			UserID:       id,
// 			AuthType:     ent.AuthTypeEmail,
// 			Identifier:   req.Email,
// 			PasswordHash: &passwordHashStr,
// 			IsVerified:   false,
// 		}
// 		if err := s.dbCredential.CreateCredential(ctx, credential); err != nil {
// 			span.AddEvent(`users.svc.create.error_creating_credential`)
// 			return err
// 		}
// 		span.AddEvent(`users.svc.create.credential_created`)

// 		// Create contact email
// 		contactEmail := &ent.ContactEmailEntity{
// 			ID:     uuid.New(),
// 			UserID: id,
// 			Email:  req.Email,
// 			IsMain: true,
// 		}
// 		if err := s.dbContactEmail.CreateContactEmail(ctx, contactEmail); err != nil {
// 			span.AddEvent(`users.svc.create.error_creating_contact_email`)
// 			return err
// 		}
// 		span.AddEvent(`users.svc.create.contact_email_created`)

// 		// Create contact phone
// 		contactPhone := &ent.ContactPhoneEntity{
// 			ID:          uuid.New(),
// 			UserID:      id,
// 			PhoneNumber: req.Phone,
// 			PhoneType:   ent.PhoneTypeMobile,
// 			CountryCode: "+66",
// 			IsMain:      true,
// 		}
// 		if err := s.dbContactPhone.CreateContactPhone(ctx, contactPhone); err != nil {
// 			span.AddEvent(`users.svc.create.error_creating_contact_phone`)
// 			return err
// 		}
// 		span.AddEvent(`users.svc.create.contact_phone_created`)

// 		return nil
// 	})

// 	if err != nil {
// 		span.AddEvent(`users.svc.create.transaction_failed`)
// 		return err
// 	}

// 	span.AddEvent(`users.svc.create.success`)
// 	return nil
// }
