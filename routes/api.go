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
		categories := system.Group("/categories")
		{
			categories.GET("/", mod.Categories.Ctl.CategoriesList)
			categories.GET("/:id", mod.Categories.Ctl.CategoriesInfo)
			categories.POST("/", mod.Categories.Ctl.CreateCategoryController)
			categories.PATCH("/:id", mod.Categories.Ctl.CategoriesUpdate)
			categories.DELETE("/:id", mod.Categories.Ctl.CategoriesDelete)
		}
		products := system.Group("/products")
		{
			products.GET("/", mod.Products.Ctl.ProductsList)
			products.GET("/:id", mod.Products.Ctl.ProductsInfo)
			products.POST("/", mod.Products.Ctl.CreateProductController)
			products.PATCH("/:id", mod.Products.Ctl.ProductsUpdate)
			products.DELETE("/:id", mod.Products.Ctl.ProductsDelete)

			products.GET("/:id/detail", mod.ProductDetails.Ctl.InfoController)
			products.POST("/:id/detail", mod.ProductDetails.Ctl.CreateController)
			products.PATCH("/:id/detail", mod.ProductDetails.Ctl.UpdateController)
			products.DELETE("/:id/detail", mod.ProductDetails.Ctl.DeleteController)

			products.GET("/:id/stock", mod.ProductStocks.Ctl.InfoController)
			products.POST("/:id/stock", mod.ProductStocks.Ctl.CreateController)
			products.PATCH("/:id/stock", mod.ProductStocks.Ctl.UpdateController)
			products.DELETE("/:id/stock", mod.ProductStocks.Ctl.DeleteController)
		}

		productStocks := system.Group("/product_stocks")
		{
			productStocks.GET("/", mod.ProductStocks.Ctl.ListController)
		}

		payments := system.Group("/payments")
		{
			payments.GET("/", mod.Payments.Ctl.PaymentsList)
			payments.GET("/:id", mod.Payments.Ctl.PaymentsInfo)
			payments.POST("/", mod.Payments.Ctl.CreatePaymentController)
			payments.PATCH("/:id", mod.Payments.Ctl.PaymentsUpdate)
			payments.DELETE("/:id", mod.Payments.Ctl.PaymentsDelete)
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
	public := r.Group("/public")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/login", mod.Auth.Ctl.LoginController)
			auth.POST("/refresh", mod.Auth.Ctl.RefreshTokenController)
		}

		members := public.Group("/members")
		{
			members.POST("/register", mod.Members.Ctl.CreateRegisterController)
		}
	}
}

func apiAuth(r *gin.RouterGroup, mod *modules.Modules) {
	auth := r.Group("/auth", mod.Auth.Ctl.AuthMiddleware())
	{
		auth.GET("/me", mod.Auth.Ctl.GetInfoController)
		auth.POST("/act-as", mod.Auth.Ctl.ActAsMemberController)
		auth.POST("/act-as/exit", mod.Auth.Ctl.ExitActAsController)

		members := auth.Group("/members")
		{
			members.GET("/", mod.Members.Ctl.ListController)
			members.GET("/:id", mod.Members.Ctl.InfoController)
			members.GET("/:id/overview", mod.Members.Ctl.MemberOverviewController)
			members.POST("/register", mod.Members.Ctl.CreateRegisterByAdminController)
			members.PATCH("/:id", mod.Members.Ctl.UpdateController)
			members.PATCH("/:id/email", mod.Members.Ctl.UpdateEmailController)
			members.PATCH("/:id/password", mod.Members.Ctl.UpdatePasswordController)
			members.DELETE("/:id", mod.Members.Ctl.DeleteController)
		}

		addresses := auth.Group("/members/:id/addresses")
		{
			addresses.GET("/", mod.Members.Ctl.ListMemberAddressesController)
			addresses.GET("/:address_id", mod.Members.Ctl.InfoMemberAddressController)
			addresses.POST("/", mod.Members.Ctl.CreateMemberAddressController)
			addresses.PATCH("/:address_id", mod.Members.Ctl.UpdateMemberAddressController)
			addresses.DELETE("/:address_id", mod.Members.Ctl.DeleteMemberAddressController)
		}

		banks := auth.Group("/members/:id/banks")
		{
			banks.GET("/", mod.Members.Ctl.ListMemberBanksController)
			banks.GET("/:bank_id", mod.Members.Ctl.InfoMemberBankController)
			banks.POST("/", mod.Members.Ctl.CreateMemberBankController)
			banks.PATCH("/:bank_id", mod.Members.Ctl.UpdateMemberBankController)
			banks.DELETE("/:bank_id", mod.Members.Ctl.DeleteMemberBankController)
		}

		files := auth.Group("/members/:id/files")
		{
			files.GET("/", mod.Members.Ctl.ListMemberFilesController)
			files.GET("/:file_row_id", mod.Members.Ctl.InfoMemberFileController)
			files.POST("/", mod.Members.Ctl.CreateMemberFileController)
			files.PATCH("/:file_row_id", mod.Members.Ctl.UpdateMemberFileController)
			files.DELETE("/:file_row_id", mod.Members.Ctl.DeleteMemberFileController)
		}

		payments := auth.Group("/members/:id/payments")
		{
			payments.GET("/", mod.Members.Ctl.ListMemberPaymentsController)
			payments.GET("/:member_payment_id", mod.Members.Ctl.InfoMemberPaymentController)
			payments.POST("/", mod.Members.Ctl.CreateMemberPaymentController)
			payments.PATCH("/:member_payment_id", mod.Members.Ctl.UpdateMemberPaymentController)
			payments.DELETE("/:member_payment_id", mod.Members.Ctl.DeleteMemberPaymentController)
		}

		wishlist := auth.Group("/members/:id/wishlist")
		{
			wishlist.GET("/", mod.Members.Ctl.ListMemberWishlistController)
			wishlist.GET("/:wishlist_id", mod.Members.Ctl.InfoMemberWishlistController)
			wishlist.POST("/", mod.Members.Ctl.CreateMemberWishlistController)
			wishlist.PATCH("/:wishlist_id", mod.Members.Ctl.UpdateMemberWishlistController)
			wishlist.DELETE("/:wishlist_id", mod.Members.Ctl.DeleteMemberWishlistController)
		}

		orders := auth.Group("/orders")
		{
			orders.GET("/", mod.Orders.Ctl.ListOrderController)
			orders.GET("/:id", mod.Orders.Ctl.InfoOrderController)
			orders.GET("/:id/timeline", mod.Orders.Ctl.TimelineOrderController)
			orders.POST("/:id/payment/confirm", mod.Orders.Ctl.ConfirmOrderPaymentController)
			orders.PATCH("/:id/payment/approve", mod.Orders.Ctl.ApproveOrderPaymentController)
			orders.PATCH("/:id/payment/reject", mod.Orders.Ctl.RejectOrderPaymentController)
			orders.POST("/", mod.Orders.Ctl.CreateOrderController)
			orders.PATCH("/:id", mod.Orders.Ctl.UpdateOrderController)
			orders.DELETE("/:id", mod.Orders.Ctl.DeleteOrderController)

			orders.GET("/:id/items", mod.Orders.Ctl.ListOrderItemController)
			orders.GET("/:id/items/:item_id", mod.Orders.Ctl.InfoOrderItemController)
			orders.POST("/:id/items", mod.Orders.Ctl.CreateOrderItemController)
			orders.PATCH("/:id/items/:item_id", mod.Orders.Ctl.UpdateOrderItemController)
			orders.DELETE("/:id/items/:item_id", mod.Orders.Ctl.DeleteOrderItemController)
		}

		carts := auth.Group("/carts")
		{
			carts.GET("/", mod.Carts.Ctl.ListCartController)
			carts.GET("/:id", mod.Carts.Ctl.InfoCartController)
			carts.POST("/", mod.Carts.Ctl.CreateCartController)
			carts.PATCH("/:id", mod.Carts.Ctl.UpdateCartController)
			carts.DELETE("/:id", mod.Carts.Ctl.DeleteCartController)

			carts.GET("/:id/items", mod.Carts.Ctl.ListCartItemController)
			carts.GET("/:id/items/:item_id", mod.Carts.Ctl.InfoCartItemController)
			carts.POST("/:id/items", mod.Carts.Ctl.CreateCartItemController)
			carts.PATCH("/:id/items/:item_id", mod.Carts.Ctl.UpdateCartItemController)
			carts.DELETE("/:id/items/:item_id", mod.Carts.Ctl.DeleteCartItemController)
		}
	}
}
