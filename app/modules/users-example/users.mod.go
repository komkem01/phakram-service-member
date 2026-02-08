package users

// import (
// 	entitiesinf "phakram-craft/app/modules/entities/inf"
// 	"phakram-craft/internal/database"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/trace"
// )

// type Module struct {
// 	Svc *Service
// 	Ctl *Controller
// }
// type (
// 	Service struct {
// 		tracer         trace.Tracer
// 		bunDB          *database.DatabaseService
// 		db             entitiesinf.UserEntity
// 		dbCredential   entitiesinf.UserCredentialEntity
// 		dbContactEmail entitiesinf.ContactEmailEntity
// 		dbContactPhone entitiesinf.ContactPhoneEntity
// 	}
// 	Controller struct {
// 		tracer trace.Tracer
// 		svc    *Service
// 	}
// )

// type Options struct {
// 	// *configDTO.Config[Config]
// 	tracer         trace.Tracer
// 	bunDB          *database.DatabaseService
// 	db             entitiesinf.UserEntity
// 	dbCredential   entitiesinf.UserCredentialEntity
// 	dbContactEmail entitiesinf.ContactEmailEntity
// 	dbContactPhone entitiesinf.ContactPhoneEntity
// }

// func New(bunDB *database.DatabaseService, db entitiesinf.UserEntity, dbCredential entitiesinf.UserCredentialEntity, dbContactEmail entitiesinf.ContactEmailEntity, dbContactPhone entitiesinf.ContactPhoneEntity) *Module {
// 	tracer := otel.Tracer("users_module")
// 	svc := newService(&Options{
// 		tracer:         tracer,
// 		bunDB:          bunDB,
// 		db:             db,
// 		dbCredential:   dbCredential,
// 		dbContactEmail: dbContactEmail,
// 		dbContactPhone: dbContactPhone,
// 	})
// 	return &Module{
// 		Svc: svc,
// 		Ctl: newController(tracer, svc),
// 	}
// }

// func newService(opt *Options) *Service {
// 	return &Service{
// 		tracer:         opt.tracer,
// 		bunDB:          opt.bunDB,
// 		db:             opt.db,
// 		dbCredential:   opt.dbCredential,
// 		dbContactEmail: opt.dbContactEmail,
// 		dbContactPhone: opt.dbContactPhone,
// 	}
// }

// func newController(trace trace.Tracer, svc *Service) *Controller {
// 	return &Controller{
// 		tracer: trace,
// 		svc:    svc,
// 	}
// }
