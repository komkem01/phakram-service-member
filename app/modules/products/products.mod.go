package products

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

type RailwayConfig struct {
	URL            string
	ServiceRoleKey string
	PublicBucket   string
	PrivateBucket  string
}

type (
	Service struct {
		tracer         trace.Tracer
		bunDB          *database.DatabaseService
		db             entitiesinf.ProductEntity
		railwayStorage *railwayStorageClient
		productFile    entitiesinf.ProductFileEntity
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer      trace.Tracer
	bunDB       *database.DatabaseService
	db          entitiesinf.ProductEntity
	railwayConf RailwayConfig
	productFile entitiesinf.ProductFileEntity
}

func New(
	bunDB *database.DatabaseService,
	db entitiesinf.ProductEntity,
	productFile entitiesinf.ProductFileEntity,
	railwayConf RailwayConfig,
) *Module {
	tracer := otel.Tracer("products_module")
	svc := newService(&Options{
		tracer:      tracer,
		bunDB:       bunDB,
		db:          db,
		productFile: productFile,
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
		db:             opt.db,
		productFile:    opt.productFile,
		railwayStorage: newRailwayStorageClient(opt.railwayConf),
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}
