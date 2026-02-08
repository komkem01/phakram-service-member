package modules

import (
	"log/slog"
	"sync"

	"phakram/app/modules/auth"
	"phakram/app/modules/banks"
	"phakram/app/modules/cart_items"
	"phakram/app/modules/carts"
	"phakram/app/modules/categories"
	"phakram/app/modules/dashboard"
	"phakram/app/modules/districts"
	"phakram/app/modules/entities"
	"phakram/app/modules/example"
	"phakram/app/modules/genders"
	"phakram/app/modules/member_accounts"
	"phakram/app/modules/member_addresses"
	"phakram/app/modules/member_banks"
	"phakram/app/modules/member_files"
	"phakram/app/modules/member_transactions"
	"phakram/app/modules/member_wishlist"
	"phakram/app/modules/members"
	"phakram/app/modules/order_items"
	"phakram/app/modules/orders"
	"phakram/app/modules/payment_files"
	"phakram/app/modules/payments"
	"phakram/app/modules/prefixes"
	"phakram/app/modules/product_details"
	"phakram/app/modules/product_files"
	"phakram/app/modules/product_stocks"
	"phakram/app/modules/products"
	"phakram/app/modules/provinces"
	"phakram/app/modules/sentry"
	"phakram/app/modules/specs"
	"phakram/app/modules/statuses"
	"phakram/app/modules/storages"
	subdistricts "phakram/app/modules/sub_districts"
	"phakram/app/modules/tiers"
	"phakram/app/modules/zipcodes"
	"phakram/internal/config"
	"phakram/internal/database"
	"phakram/internal/log"
	"phakram/internal/otel/collector"

	exampletwo "phakram/app/modules/example-two"

	appConf "phakram/config"
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
	Genders            *genders.Module
	Prefixes           *prefixes.Module
	Banks              *banks.Module
	Storages           *storages.Module
	Members            *members.Module
	MemberAccounts     *member_accounts.Module
	MemberAddresses    *member_addresses.Module
	MemberBanks        *member_banks.Module
	MemberFiles        *member_files.Module
	MemberTransactions *member_transactions.Module
	MemberWishlist     *member_wishlist.Module
	Dashboard          *dashboard.Module
	Orders             *orders.Module
	OrderItems         *order_items.Module
	Carts              *carts.Module
	CartItems          *cart_items.Module
	Payments           *payments.Module
	PaymentFiles       *payment_files.Module
	Products           *products.Module
	ProductDetails     *product_details.Module
	ProductStocks      *product_stocks.Module
	ProductFiles       *product_files.Module
	Categories         *categories.Module
	Provinces          *provinces.Module
	Districts          *districts.Module
	SubDistricts       *subdistricts.Module
	Zipcodes           *zipcodes.Module
	Statuses           *statuses.Module
	Tiers              *tiers.Module
	Auth               *auth.Module
}

func modulesInit() {
	confMod := config.New(&appConf.App)
	specsMod := specs.New(config.Conf[specs.Config](confMod.Svc))
	conf := confMod.Svc.Config()

	logMod := log.New(config.Conf[log.Option](confMod.Svc))
	otel := collector.New(config.Conf[collector.Config](confMod.Svc))
	log := log.With(slog.String("module", "modules"))

	sentryMod := sentry.New(config.Conf[sentry.Config](confMod.Svc))

	db := database.New(conf.Database.Sql)
	entitiesMod := entities.New(db.Svc.DB())
	exampleMod := example.New(config.Conf[example.Config](confMod.Svc), entitiesMod.Svc)
	exampleMod2 := exampletwo.New(config.Conf[exampletwo.Config](confMod.Svc), entitiesMod.Svc)
	// kafka := kafka.New(&conf.Kafka)
	gendersMod := genders.New(db.Svc, entitiesMod.Svc)
	prefixesMod := prefixes.New(db.Svc, entitiesMod.Svc)
	banksMod := banks.New(db.Svc, entitiesMod.Svc)
	storagesMod := storages.New(db.Svc, entitiesMod.Svc)
	membersMod := members.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	memberAccountsMod := member_accounts.New(db.Svc, entitiesMod.Svc)
	memberAddressesMod := member_addresses.New(db.Svc, entitiesMod.Svc)
	memberBanksMod := member_banks.New(db.Svc, entitiesMod.Svc)
	memberFilesMod := member_files.New(db.Svc, entitiesMod.Svc)
	memberTransactionsMod := member_transactions.New(db.Svc, entitiesMod.Svc)
	memberWishlistMod := member_wishlist.New(db.Svc, entitiesMod.Svc)
	dashboardMod := dashboard.New(db.Svc)
	ordersMod := orders.New(db.Svc, entitiesMod.Svc)
	orderItemsMod := order_items.New(db.Svc, entitiesMod.Svc)
	cartsMod := carts.New(db.Svc, entitiesMod.Svc)
	cartItemsMod := cart_items.New(db.Svc, entitiesMod.Svc)
	paymentsMod := payments.New(db.Svc, entitiesMod.Svc)
	paymentFilesMod := payment_files.New(db.Svc, entitiesMod.Svc)
	productsMod := products.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	productDetailsMod := product_details.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	productStocksMod := product_stocks.New(db.Svc, entitiesMod.Svc)
	productFilesMod := product_files.New(db.Svc, entitiesMod.Svc)
	categoriesMod := categories.New(db.Svc, entitiesMod.Svc)
	provincesMod := provinces.New(db.Svc, entitiesMod.Svc)
	districtsMod := districts.New(db.Svc, entitiesMod.Svc)
	subDistrictsMod := subdistricts.New(db.Svc, entitiesMod.Svc)
	zipcodesMod := zipcodes.New(db.Svc, entitiesMod.Svc)
	statusesMod := statuses.New(db.Svc, entitiesMod.Svc)
	tiersMod := tiers.New(db.Svc, entitiesMod.Svc)
	authMod := auth.New(db.Svc, entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc, conf.AppKey)
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
		Genders:            gendersMod,
		Prefixes:           prefixesMod,
		Banks:              banksMod,
		Storages:           storagesMod,
		Members:            membersMod,
		MemberAccounts:     memberAccountsMod,
		MemberAddresses:    memberAddressesMod,
		MemberBanks:        memberBanksMod,
		MemberFiles:        memberFilesMod,
		MemberTransactions: memberTransactionsMod,
		MemberWishlist:     memberWishlistMod,
		Dashboard:          dashboardMod,
		Orders:             ordersMod,
		OrderItems:         orderItemsMod,
		Carts:              cartsMod,
		CartItems:          cartItemsMod,
		Payments:           paymentsMod,
		PaymentFiles:       paymentFilesMod,
		Products:           productsMod,
		ProductDetails:     productDetailsMod,
		ProductStocks:      productStocksMod,
		ProductFiles:       productFilesMod,
		Categories:         categoriesMod,
		Provinces:          provincesMod,
		Districts:          districtsMod,
		SubDistricts:       subDistrictsMod,
		Zipcodes:           zipcodesMod,
		Statuses:           statusesMod,
		Tiers:              tiersMod,
		Auth:               authMod,
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
