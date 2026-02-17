package productstocks

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

func New(bunDB *database.DatabaseService) *Module {
	tracer := otel.Tracer("product_stocks_module")
	svc := &Service{tracer: tracer, bunDB: bunDB}
	return &Module{Svc: svc, Ctl: &Controller{tracer: tracer, svc: svc}}
}
