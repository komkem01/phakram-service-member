package reviews

import (
	"phakram/internal/database"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	Svc *Service
	Ctl *Controller
}

type RailwayConfig struct {
	URL            string
	ServiceRoleKey string
	PublicBucket   string
	ReviewBucket   string
	PrivateBucket  string
}

type (
	Service struct {
		tracer         trace.Tracer
		bunDB          *database.DatabaseService
		railwayStorage *railwayStorageClient
	}

	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer      trace.Tracer
	bunDB       *database.DatabaseService
	railwayConf RailwayConfig
}

func New(bunDB *database.DatabaseService, railwayConf RailwayConfig) *Module {
	tracer := otel.Tracer("reviews_module")
	svc := newService(&Options{
		tracer:      tracer,
		bunDB:       bunDB,
		railwayConf: railwayConf,
	})

	return &Module{
		Svc: svc,
		Ctl: newController(tracer, svc),
	}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:         opt.tracer,
		bunDB:          opt.bunDB,
		railwayStorage: newRailwayStorageClient(opt.railwayConf),
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}
