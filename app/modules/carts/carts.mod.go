package carts

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
		cart   entitiesinf.CartEntity
		item   entitiesinf.CartItemEntity
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer trace.Tracer
	bunDB  *database.DatabaseService
	cart   entitiesinf.CartEntity
	item   entitiesinf.CartItemEntity
}

func New(bunDB *database.DatabaseService, cart entitiesinf.CartEntity, item entitiesinf.CartItemEntity) *Module {
	tracer := otel.Tracer("carts_module")
	svc := newService(&Options{tracer: tracer, bunDB: bunDB, cart: cart, item: item})
	return &Module{Svc: svc, Ctl: newController(tracer, svc)}
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, bunDB: opt.bunDB, cart: opt.cart, item: opt.item}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}
