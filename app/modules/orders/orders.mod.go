package orders

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
		tracer trace.Tracer
		bunDB  *database.DatabaseService
		order  entitiesinf.OrderEntity
		item   entitiesinf.OrderItemEntity
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer trace.Tracer
	bunDB  *database.DatabaseService
	order  entitiesinf.OrderEntity
	item   entitiesinf.OrderItemEntity
}

func New(bunDB *database.DatabaseService, order entitiesinf.OrderEntity, item entitiesinf.OrderItemEntity) *Module {
	tracer := otel.Tracer("orders_module")
	svc := newService(&Options{tracer: tracer, bunDB: bunDB, order: order, item: item})
	return &Module{Svc: svc, Ctl: newController(tracer, svc)}
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, bunDB: opt.bunDB, order: opt.order, item: opt.item}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}
