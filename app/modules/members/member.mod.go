package members

import (
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/internal/database"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	Svc *Service
	Ctl *Controller
}

type (
	Service struct {
		tracer    trace.Tracer
		bunDB     *database.DatabaseService
		db        entitiesinf.MemberEntity
		dbAccount entitiesinf.MemberAccountEntity
		audit     entitiesinf.AuditLogEntity
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer    trace.Tracer
	bunDB     *database.DatabaseService
	db        entitiesinf.MemberEntity
	dbAccount entitiesinf.MemberAccountEntity
	audit     entitiesinf.AuditLogEntity
}

func New(bunDB *database.DatabaseService, db entitiesinf.MemberEntity, dbAccount entitiesinf.MemberAccountEntity, audit entitiesinf.AuditLogEntity) *Module {
	tracer := otel.Tracer("members_module")
	svc := newService(&Options{
		tracer:    tracer,
		bunDB:     bunDB,
		db:        db,
		dbAccount: dbAccount,
		audit:     audit,
	})
	return &Module{
		Svc: svc,
		Ctl: newController(tracer, svc),
	}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:    opt.tracer,
		bunDB:     opt.bunDB,
		db:        opt.db,
		dbAccount: opt.dbAccount,
		audit:     opt.audit,
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}
