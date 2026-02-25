package modules

import (
	"log/slog"
	"phakram/app/modules/auth"
	"phakram/app/modules/banks"
	"phakram/app/modules/carts"
	"phakram/app/modules/categories"
	"phakram/app/modules/contact"
	"phakram/app/modules/districts"
	"phakram/app/modules/entities"
	"phakram/app/modules/example"
	exampletwo "phakram/app/modules/example-two"
	"phakram/app/modules/genders"
	"phakram/app/modules/members"
	"phakram/app/modules/orders"
	"phakram/app/modules/payments"
	"phakram/app/modules/prefixes"
	productdetails "phakram/app/modules/product_details"
	productstocks "phakram/app/modules/product_stocks"
	"phakram/app/modules/products"
	"phakram/app/modules/promotions"
	"phakram/app/modules/provinces"
	"phakram/app/modules/reviews"
	"phakram/app/modules/sentry"
	"phakram/app/modules/specs"
	"phakram/app/modules/statuses"
	"phakram/app/modules/storages"
	subdistricts "phakram/app/modules/sub_districts"
	systembankaccounts "phakram/app/modules/system_bank_accounts"
	"phakram/app/modules/tiers"
	"phakram/app/modules/zipcodes"
	appConf "phakram/config"
	"phakram/internal/config"
	"phakram/internal/database"
	"phakram/internal/log"
	"phakram/internal/otel/collector"
	"strings"
	"sync"
	// "phakram/app/modules/kafka"
)

type Modules struct {
	Conf   *config.Module[appConf.Config]
	Specs  *specs.Module
	Log    *log.Module
	OTEL   *collector.Module
	Sentry *sentry.Module
	DB     *database.DatabaseModule
	ENT    *entities.Module
	// Kafka *kafka.Module
	Example            *example.Module
	Example2           *exampletwo.Module
	Contact            *contact.Module
	Genders            *genders.Module
	Prefixes           *prefixes.Module
	Banks              *banks.Module
	Provinces          *provinces.Module
	Districts          *districts.Module
	SubDistricts       *subdistricts.Module
	Zipcodes           *zipcodes.Module
	Statuses           *statuses.Module
	Tiers              *tiers.Module
	SystemBankAccounts *systembankaccounts.Module
	Categories         *categories.Module
	Products           *products.Module
	ProductDetails     *productdetails.Module
	ProductStocks      *productstocks.Module
	Storages           *storages.Module
	Auth               *auth.Module
	Members            *members.Module
	Orders             *orders.Module
	Payments           *payments.Module
	Carts              *carts.Module
	Promotions         *promotions.Module
	Reviews            *reviews.Module
}

func modulesInit() {
	confMod := config.New(&appConf.App)
	specsMod := specs.New(config.Conf[specs.Config](confMod.Svc))
	conf := confMod.Svc.Config()

	logMod := log.New(config.Conf[log.Option](confMod.Svc))
	otel := collector.New(config.Conf[collector.Config](confMod.Svc))
	log := log.With(slog.String("module", "modules"))
	log.With(
		slog.String("railway_storage_url", strings.TrimSpace(conf.RailwayStorage.URL)),
		slog.String("railway_public_bucket", strings.TrimSpace(conf.RailwayStorage.PublicBucket)),
		slog.String("railway_review_bucket", strings.TrimSpace(conf.RailwayStorage.ReviewBucket)),
		slog.Bool("railway_secret_access_key_set", strings.TrimSpace(conf.RailwayStorage.ServiceRoleKey) != ""),
	).Infof("railway storage config loaded")

	sentryMod := sentry.New(config.Conf[sentry.Config](confMod.Svc))

	db := database.New(conf.Database.Sql)
	entitiesMod := entities.New(db.Svc.DB())
	exampleMod := example.New(config.Conf[example.Config](confMod.Svc), entitiesMod.Svc)
	exampleMod2 := exampletwo.New(config.Conf[exampletwo.Config](confMod.Svc), entitiesMod.Svc)
	// kafka := kafka.New(&conf.Kafka)
	gendersMod := genders.New(db.Svc, entitiesMod.Svc)
	prefixesMod := prefixes.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc)
	banksMod := banks.New(db.Svc, entitiesMod.Svc)
	provincesMod := provinces.New(db.Svc, entitiesMod.Svc)
	districtsMod := districts.New(db.Svc, entitiesMod.Svc)
	subDistrictsMod := subdistricts.New(db.Svc, entitiesMod.Svc)
	zipcodesMod := zipcodes.New(db.Svc, entitiesMod.Svc)
	statusesMod := statuses.New(db.Svc, entitiesMod.Svc)
	tiersMod := tiers.New(db.Svc, entitiesMod.Svc)
	systemBankAccountsMod := systembankaccounts.New(db.Svc)
	categoriesMod := categories.New(db.Svc, entitiesMod.Svc)
	productsMod := products.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc, products.RailwayConfig{
		URL:            conf.RailwayStorage.URL,
		ServiceRoleKey: conf.RailwayStorage.ServiceRoleKey,
		PublicBucket:   conf.RailwayStorage.PublicBucket,
		PrivateBucket:  conf.RailwayStorage.PrivateBucket,
	})
	productDetailsMod := productdetails.New(db.Svc)
	productStocksMod := productstocks.New(db.Svc)
	storagesMod := storages.New(db.Svc, entitiesMod.Svc, storages.RailwayConfig{
		URL:            conf.RailwayStorage.URL,
		ServiceRoleKey: conf.RailwayStorage.ServiceRoleKey,
		PublicBucket:   conf.RailwayStorage.PublicBucket,
		PrivateBucket:  conf.RailwayStorage.PrivateBucket,
	})
	authMod := auth.New(db.Svc, conf.AppKey)
	membersMod := members.New(
		db.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
	)
	ordersMod := orders.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc, orders.RailwayConfig{
		URL:            conf.RailwayStorage.URL,
		ServiceRoleKey: conf.RailwayStorage.ServiceRoleKey,
		PublicBucket:   conf.RailwayStorage.PublicBucket,
		PrivateBucket:  conf.RailwayStorage.PrivateBucket,
	})
	contactMod := contact.New(db.Svc, &conf.Contact)
	paymentsMod := payments.New(db.Svc, entitiesMod.Svc)
	cartsMod := carts.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc)
	promotionsMod := promotions.New(db.Svc)
	reviewsMod := reviews.New(db.Svc, reviews.RailwayConfig{
		URL:            conf.RailwayStorage.URL,
		ServiceRoleKey: conf.RailwayStorage.ServiceRoleKey,
		PublicBucket:   conf.RailwayStorage.PublicBucket,
		ReviewBucket:   conf.RailwayStorage.ReviewBucket,
		PrivateBucket:  conf.RailwayStorage.PrivateBucket,
	})
	mod = &Modules{
		Conf:               confMod,
		Specs:              specsMod,
		Log:                logMod,
		OTEL:               otel,
		Sentry:             sentryMod,
		DB:                 db,
		ENT:                entitiesMod,
		Example:            exampleMod,
		Example2:           exampleMod2,
		Contact:            contactMod,
		Genders:            gendersMod,
		Prefixes:           prefixesMod,
		Banks:              banksMod,
		Provinces:          provincesMod,
		Districts:          districtsMod,
		SubDistricts:       subDistrictsMod,
		Zipcodes:           zipcodesMod,
		Statuses:           statusesMod,
		Tiers:              tiersMod,
		SystemBankAccounts: systemBankAccountsMod,
		Categories:         categoriesMod,
		Products:           productsMod,
		ProductDetails:     productDetailsMod,
		ProductStocks:      productStocksMod,
		Storages:           storagesMod,
		Auth:               authMod,
		Members:            membersMod,
		Orders:             ordersMod,
		Payments:           paymentsMod,
		Carts:              cartsMod,
		Promotions:         promotionsMod,
		Reviews:            reviewsMod,
	}

	log.Infof("all modules initialized")
}

var (
	once sync.Once
	mod  *Modules
)

func Get() *Modules {
	once.Do(modulesInit)

	return mod
}
