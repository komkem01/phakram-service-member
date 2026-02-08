package auth

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
		tracer       trace.Tracer
		bunDB        *database.DatabaseService
		db           entitiesinf.MemberEntity
		dbAccount    entitiesinf.MemberAccountEntity
		dbTrasaction entitiesinf.MemberTransactionEntity
		secret       string
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer       trace.Tracer
	bunDB        *database.DatabaseService
	db           entitiesinf.MemberEntity
	dbAccount    entitiesinf.MemberAccountEntity
	dbTrasaction entitiesinf.MemberTransactionEntity
	secret       string
}

func New(bunDB *database.DatabaseService, db entitiesinf.MemberEntity, dbAccount entitiesinf.MemberAccountEntity, dbTrasaction entitiesinf.MemberTransactionEntity, secret string) *Module {
	tracer := otel.Tracer("auth_module")
	svc := newService(&Options{
		tracer:       tracer,
		bunDB:        bunDB,
		db:           db,
		dbAccount:    dbAccount,
		dbTrasaction: dbTrasaction,
		secret:       secret,
	})
	return &Module{
		Svc: svc,
		Ctl: newController(tracer, svc),
	}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:       opt.tracer,
		bunDB:        opt.bunDB,
		db:           opt.db,
		dbAccount:    opt.dbAccount,
		dbTrasaction: opt.dbTrasaction,
		secret:       opt.secret,
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}
