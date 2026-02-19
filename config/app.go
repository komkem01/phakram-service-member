package config

import (
	"phakram/app/modules/example"
	exampletwo "phakram/app/modules/example-two"
	"phakram/app/modules/sentry"
	"phakram/app/modules/specs"
	"phakram/internal/kafka"
	"phakram/internal/log"
	"phakram/internal/otel/collector"
)

// Config is a struct that contains all the configuration of the application.
type Config struct {
	Database Database
	Supabase SupabaseConfig

	AppName     string
	AppKey      string
	Environment string
	Specs       specs.Config
	Debug       bool

	Port           int
	HttpJsonNaming string

	SslCaPath      string
	SslPrivatePath string
	SslCertPath    string

	Otel   collector.Config
	Sentry sentry.Config

	Kafka kafka.Config
	Log   log.Option

	Example example.Config

	ExampleTwo exampletwo.Config
}

type SupabaseConfig struct {
	URL            string
	ServiceRoleKey string
	PublicBucket   string
	PrivateBucket  string
}

var App = Config{
	Specs: specs.Config{
		Version: "v1",
	},
	Database: database,
	Kafka:    kafkaConf,
	Supabase: SupabaseConfig{},

	AppName: "go_app",
	Port:    8081,
	AppKey:  "secret",
	Debug:   false,

	HttpJsonNaming: "snake_case",

	SslCaPath:      "phakram/cert/ca.pem",
	SslPrivatePath: "phakram/cert/server.pem",
	SslCertPath:    "phakram/cert/server-key.pem",

	Otel: collector.Config{
		CollectorEndpoint: "",
		LogMode:           "noop",
		TraceMode:         "noop",
		MetricMode:        "noop",
		TraceRatio:        0.01,
	},
}
