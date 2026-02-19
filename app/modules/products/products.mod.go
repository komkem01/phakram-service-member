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

type SupabaseConfig struct {
	URL            string
	ServiceRoleKey string
	PublicBucket   string
	PrivateBucket  string
}

type (
	Service struct {
		tracer      trace.Tracer
		bunDB       *database.DatabaseService
		db          entitiesinf.ProductEntity
		supabase    *supabaseStorageClient
		productFile entitiesinf.ProductFileEntity
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer       trace.Tracer
	bunDB        *database.DatabaseService
	db           entitiesinf.ProductEntity
	supabaseConf SupabaseConfig
	productFile  entitiesinf.ProductFileEntity
}

func New(
	bunDB *database.DatabaseService,
	db entitiesinf.ProductEntity,
	productFile entitiesinf.ProductFileEntity,
	supabaseConf SupabaseConfig,
) *Module {
	tracer := otel.Tracer("products_module")
	svc := newService(&Options{
		tracer:       tracer,
		bunDB:        bunDB,
		db:           db,
		productFile:  productFile,
		supabaseConf: supabaseConf,
	})
	return &Module{
		Svc: svc,
		Ctl: newController(tracer, svc),
	}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:      opt.tracer,
		bunDB:       opt.bunDB,
		db:          opt.db,
		productFile: opt.productFile,
		supabase:    newSupabaseStorageClient(opt.supabaseConf),
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}
