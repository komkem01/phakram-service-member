package contact

import (
	"phakram/internal/database"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Config struct {
	RecipientEmail string
	Mail           MailConfig
}

type Module struct {
	Svc *Service
	Ctl *Controller
}

type (
	Service struct {
		tracer trace.Tracer
		bunDB  *database.DatabaseService
		conf   *Config
	}

	Controller struct {
		tracer trace.Tracer
		svc    *Service
	}
)

type Options struct {
	tracer trace.Tracer
	bunDB  *database.DatabaseService
	conf   *Config
}

func New(bunDB *database.DatabaseService, conf *Config) *Module {
	tracer := otel.Tracer("contact_module")
	svc := newService(&Options{
		tracer: tracer,
		bunDB:  bunDB,
		conf:   conf,
	})

	return &Module{
		Svc: svc,
		Ctl: newController(tracer, svc),
	}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		bunDB:  opt.bunDB,
		conf:   opt.conf,
	}
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}
