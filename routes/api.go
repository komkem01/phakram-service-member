package routes

import (
	"fmt"
	"net/http"

	"phakram/app/modules"

	"github.com/gin-gonic/gin"
)

func WarpH(router *gin.RouterGroup, prefix string, handler http.Handler) {
	router.Any(fmt.Sprintf("%s/*w", prefix), gin.WrapH(http.StripPrefix(fmt.Sprintf("%s%s", router.BasePath(), prefix), handler)))
}

func api(r *gin.RouterGroup, mod *modules.Modules) {
	r.GET("/example/:id", mod.Example.Ctl.Get)
	r.GET("/example-http", mod.Example.Ctl.GetHttpReq)
	r.POST("/example", mod.Example.Ctl.Create)
}

func apiSystem(r *gin.RouterGroup, mod *modules.Modules) {
	// Public routes (no authentication required)
	system := r.Group("/system")
	{
		genders := system.Group("/genders")
		{
			genders.GET("/", mod.Genders.Ctl.GendersList)
			genders.GET("/:id", mod.Genders.Ctl.GendersInfo)
			genders.POST("/", mod.Genders.Ctl.CreateGenderController)
			genders.PATCH("/:id", mod.Genders.Ctl.GendersUpdate)
			genders.DELETE("/:id", mod.Genders.Ctl.GendersDelete)
		}
		prefixes := system.Group("/prefixes")
		{
			prefixes.GET("/", mod.Prefixes.Ctl.PrefixesList)
			prefixes.GET("/:id", mod.Prefixes.Ctl.PrefixesInfo)
			prefixes.POST("/", mod.Prefixes.Ctl.CreatePrefixController)
			prefixes.PATCH("/:id", mod.Prefixes.Ctl.PrefixesUpdate)
			prefixes.DELETE("/:id", mod.Prefixes.Ctl.PrefixesDelete)
		}
		banks := system.Group("/banks")
		{
			banks.GET("/", mod.Banks.Ctl.BanksList)
			banks.GET("/:id", mod.Banks.Ctl.BanksInfo)
			banks.POST("/", mod.Banks.Ctl.CreateBankController)
			banks.PATCH("/:id", mod.Banks.Ctl.BanksUpdate)
			banks.DELETE("/:id", mod.Banks.Ctl.BanksDelete)
		}
		provinces := system.Group("/provinces")
		{
			provinces.GET("/", mod.Provinces.Ctl.ProvincesList)
			provinces.GET("/:id", mod.Provinces.Ctl.ProvincesInfo)
			provinces.POST("/", mod.Provinces.Ctl.CreateProvinceController)
			provinces.PATCH("/:id", mod.Provinces.Ctl.ProvincesUpdate)
			provinces.DELETE("/:id", mod.Provinces.Ctl.ProvincesDelete)
		}
		districts := system.Group("/districts")
		{
			districts.GET("/", mod.Districts.Ctl.DistrictsList)
			districts.GET("/:id", mod.Districts.Ctl.DistrictsInfo)
			districts.POST("/", mod.Districts.Ctl.CreateDistrictController)
			districts.PATCH("/:id", mod.Districts.Ctl.DistrictsUpdate)
			districts.DELETE("/:id", mod.Districts.Ctl.DistrictsDelete)
		}
		subDistricts := system.Group("/sub_districts")
		{
			subDistricts.GET("/", mod.SubDistricts.Ctl.SubDistrictsList)
			subDistricts.GET("/:id", mod.SubDistricts.Ctl.SubDistrictsInfo)
			subDistricts.POST("/", mod.SubDistricts.Ctl.CreateSubDistrictController)
			subDistricts.PATCH("/:id", mod.SubDistricts.Ctl.SubDistrictsUpdate)
			subDistricts.DELETE("/:id", mod.SubDistricts.Ctl.SubDistrictsDelete)
		}
		zipcodes := system.Group("/zipcodes")
		{
			zipcodes.GET("/", mod.Zipcodes.Ctl.ZipcodesList)
			zipcodes.GET("/:id", mod.Zipcodes.Ctl.ZipcodesInfo)
			zipcodes.POST("/", mod.Zipcodes.Ctl.CreateZipcodeController)
			zipcodes.PATCH("/:id", mod.Zipcodes.Ctl.ZipcodesUpdate)
			zipcodes.DELETE("/:id", mod.Zipcodes.Ctl.ZipcodesDelete)
		}
		statuses := system.Group("/statuses")
		{
			statuses.GET("/", mod.Statuses.Ctl.StatusesList)
			statuses.GET("/:id", mod.Statuses.Ctl.StatusesInfo)
			statuses.POST("/", mod.Statuses.Ctl.CreateStatusController)
			statuses.PATCH("/:id", mod.Statuses.Ctl.StatusesUpdate)
			statuses.DELETE("/:id", mod.Statuses.Ctl.StatusesDelete)
		}
		tiers := system.Group("/tiers")
		{
			tiers.GET("/", mod.Tiers.Ctl.TiersList)
			tiers.GET("/:id", mod.Tiers.Ctl.TiersInfo)
			tiers.POST("/", mod.Tiers.Ctl.CreateTierController)
			tiers.PATCH("/:id", mod.Tiers.Ctl.TiersUpdate)
			tiers.DELETE("/:id", mod.Tiers.Ctl.TiersDelete)
		}
	}
}

func apiStorage(r *gin.RouterGroup, mod *modules.Modules) {
	auth := r.Group("/auth", mod.Auth.Ctl.AuthMiddleware())
	{
		storages := auth.Group("/storages")
		{
			storages.GET("/", mod.Storages.Ctl.StoragesList)
			storages.GET("/:id", mod.Storages.Ctl.StoragesInfo)
			storages.POST("/", mod.Storages.Ctl.CreateStorageController)
			storages.PATCH("/:id", mod.Storages.Ctl.StoragesUpdate)
			storages.DELETE("/:id", mod.Storages.Ctl.StoragesDelete)
		}
	}
}

func apiPublic(r *gin.RouterGroup, mod *modules.Modules) {
	Public := r.Group("/public")
	{
		auth := Public.Group("/auth")
		{
			auth.POST("/register", mod.Auth.Ctl.RegisterMemberController)
			auth.POST("/login", mod.Auth.Ctl.LoginMemberController)
			auth.POST("/refresh", mod.Auth.Ctl.RefreshTokenController)
			auth.GET("/me", mod.Auth.Ctl.MeController)
		}
	}
}
