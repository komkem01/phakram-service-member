package prefixes

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
		tracer   trace.Tracer
		bunDB    *database.DatabaseService
		db       entitiesinf.PrefixEntity
		dbGender entitiesinf.GenderEntity
	}
	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	// *configDTO.Config[Config]
	tracer   trace.Tracer
	bunDB    *database.DatabaseService
	db       entitiesinf.PrefixEntity
	dbGender entitiesinf.GenderEntity
}

func New(bunDB *database.DatabaseService, db entitiesinf.PrefixEntity, dbGender entitiesinf.GenderEntity) *Module {
	tracer := otel.Tracer("prefixes_module")
	svc := newService(&Options{
		tracer:   tracer,
		bunDB:    bunDB,
		db:       db,
		dbGender: dbGender,
	})
	return &Module{
		Svc: svc,
		Ctl: newController(tracer, svc),
	}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:   opt.tracer,
		bunDB:    opt.bunDB,
		db:       opt.db,
		dbGender: opt.dbGender,
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}
