package storages

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
		tracer   trace.Tracer
		bunDB    *database.DatabaseService
		db       entitiesinf.StorageEntity
		supabase *supabaseStorageClient
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer       trace.Tracer
	bunDB        *database.DatabaseService
	db           entitiesinf.StorageEntity
	supabaseConf SupabaseConfig
}

func New(bunDB *database.DatabaseService, db entitiesinf.StorageEntity, supabaseConf SupabaseConfig) *Module {
	tracer := otel.Tracer("storages_module")
	svc := newService(&Options{
		tracer:       tracer,
		bunDB:        bunDB,
		db:           db,
		supabaseConf: supabaseConf,
	})
	return &Module{Svc: svc, Ctl: newController(tracer, svc)}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:   opt.tracer,
		bunDB:    opt.bunDB,
		db:       opt.db,
		supabase: newSupabaseStorageClient(opt.supabaseConf),
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}
