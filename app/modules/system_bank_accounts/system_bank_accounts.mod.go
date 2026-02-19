package systembankaccounts

import (
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
		tracer trace.Tracer
		bunDB  *database.DatabaseService
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer trace.Tracer
	bunDB  *database.DatabaseService
}

func New(bunDB *database.DatabaseService) *Module {
	tracer := otel.Tracer("system_bank_accounts_module")
	svc := newService(&Options{tracer: tracer, bunDB: bunDB})
	return &Module{Svc: svc, Ctl: newController(tracer, svc)}
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, bunDB: opt.bunDB}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}
