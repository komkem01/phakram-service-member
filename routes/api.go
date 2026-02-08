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

func apiMember(r *gin.RouterGroup, mod *modules.Modules) {
	auth := r.Group("/auth", mod.Auth.Ctl.AuthMiddleware())
	{
		members := auth.Group("/members")
		{
			members.GET("/", mod.Members.Ctl.MembersList)
			members.GET("/:id", mod.Members.Ctl.MembersInfo)
			members.POST("/", mod.Members.Ctl.CreateMemberController)
			members.PATCH("/:id", mod.Members.Ctl.MembersUpdate)
			members.DELETE("/:id", mod.Members.Ctl.MembersDelete)
		}
		memberAccounts := auth.Group("/member_accounts")
		{
			memberAccounts.GET("/", mod.MemberAccounts.Ctl.MemberAccountsList)
			memberAccounts.GET("/by_email", mod.MemberAccounts.Ctl.MemberAccountsInfoByEmail)
			memberAccounts.GET("/:id", mod.MemberAccounts.Ctl.MemberAccountsInfo)
			memberAccounts.POST("/", mod.MemberAccounts.Ctl.CreateMemberAccountController)
			memberAccounts.PATCH("/:id", mod.MemberAccounts.Ctl.MemberAccountsUpdate)
			memberAccounts.DELETE("/:id", mod.MemberAccounts.Ctl.MemberAccountsDelete)
		}
		memberAddresses := auth.Group("/member_addresses")
		{
			memberAddresses.GET("/", mod.MemberAddresses.Ctl.MemberAddressesList)
			memberAddresses.GET("/:id", mod.MemberAddresses.Ctl.MemberAddressesInfo)
			memberAddresses.POST("/", mod.MemberAddresses.Ctl.CreateMemberAddressController)
			memberAddresses.PATCH("/:id", mod.MemberAddresses.Ctl.MemberAddressesUpdate)
			memberAddresses.DELETE("/:id", mod.MemberAddresses.Ctl.MemberAddressesDelete)
		}
		memberBanks := auth.Group("/member_banks")
		{
			memberBanks.GET("/", mod.MemberBanks.Ctl.MemberBanksList)
			memberBanks.GET("/:id", mod.MemberBanks.Ctl.MemberBanksInfo)
			memberBanks.POST("/", mod.MemberBanks.Ctl.CreateMemberBankController)
			memberBanks.PATCH("/:id", mod.MemberBanks.Ctl.MemberBanksUpdate)
			memberBanks.DELETE("/:id", mod.MemberBanks.Ctl.MemberBanksDelete)
		}
		memberFiles := auth.Group("/member_files")
		{
			memberFiles.GET("/", mod.MemberFiles.Ctl.MemberFilesList)
			memberFiles.GET("/:id", mod.MemberFiles.Ctl.MemberFilesInfo)
			memberFiles.POST("/", mod.MemberFiles.Ctl.CreateMemberFileController)
			memberFiles.PATCH("/:id", mod.MemberFiles.Ctl.MemberFilesUpdate)
			memberFiles.DELETE("/:id", mod.MemberFiles.Ctl.MemberFilesDelete)
		}
		memberTransactions := auth.Group("/member_transactions")
		{
			memberTransactions.GET("/", mod.MemberTransactions.Ctl.MemberTransactionsList)
			memberTransactions.GET("/:id", mod.MemberTransactions.Ctl.MemberTransactionsInfo)
			memberTransactions.POST("/", mod.MemberTransactions.Ctl.CreateMemberTransactionController)
			memberTransactions.PATCH("/:id", mod.MemberTransactions.Ctl.MemberTransactionsUpdate)
			memberTransactions.DELETE("/:id", mod.MemberTransactions.Ctl.MemberTransactionsDelete)
		}
		memberWishlist := auth.Group("/member_wishlist")
		{
			memberWishlist.GET("/", mod.MemberWishlist.Ctl.MemberWishlistList)
			memberWishlist.GET("/:id", mod.MemberWishlist.Ctl.MemberWishlistInfo)
			memberWishlist.POST("/", mod.MemberWishlist.Ctl.CreateMemberWishlistController)
			memberWishlist.PATCH("/:id", mod.MemberWishlist.Ctl.MemberWishlistUpdate)
			memberWishlist.DELETE("/:id", mod.MemberWishlist.Ctl.MemberWishlistDelete)
		}
		dashboard := auth.Group("/dashboard")
		{
			dashboard.GET("/", mod.Dashboard.Ctl.Summary)
		}
	}
}

func apiProduct(r *gin.RouterGroup, mod *modules.Modules) {
	auth := r.Group("/auth", mod.Auth.Ctl.AuthMiddleware())
	{
		products := auth.Group("/products")
		{
			products.GET("/", mod.Products.Ctl.ProductsList)
			products.GET("/:id", mod.Products.Ctl.ProductsInfo)
			products.POST("/", mod.Products.Ctl.CreateProductController)
			products.PATCH("/:id", mod.Products.Ctl.ProductsUpdate)
			products.DELETE("/:id", mod.Products.Ctl.ProductsDelete)
		}
		product_details := auth.Group("/product_details")
		{
			product_details.GET("/", mod.ProductDetails.Ctl.ProductDetailsList)
			product_details.GET("/:id", mod.ProductDetails.Ctl.ProductDetailsInfo)
			product_details.POST("/", mod.ProductDetails.Ctl.CreateProductDetailController)
			product_details.PATCH("/:id", mod.ProductDetails.Ctl.ProductDetailsUpdate)
			product_details.DELETE("/:id", mod.ProductDetails.Ctl.ProductDetailsDelete)
		}
		product_stocks := auth.Group("/product_stocks")
		{
			product_stocks.GET("/", mod.ProductStocks.Ctl.ProductStocksList)
			product_stocks.GET("/:id", mod.ProductStocks.Ctl.ProductStocksInfo)
			product_stocks.POST("/", mod.ProductStocks.Ctl.CreateProductStockController)
			product_stocks.PATCH("/:id", mod.ProductStocks.Ctl.ProductStocksUpdate)
			product_stocks.DELETE("/:id", mod.ProductStocks.Ctl.ProductStocksDelete)
		}
		product_files := auth.Group("/product_files")
		{
			product_files.GET("/", mod.ProductFiles.Ctl.ProductFilesList)
			product_files.GET("/:id", mod.ProductFiles.Ctl.ProductFilesInfo)
			product_files.POST("/", mod.ProductFiles.Ctl.CreateProductFileController)
			product_files.PATCH("/:id", mod.ProductFiles.Ctl.ProductFilesUpdate)
			product_files.DELETE("/:id", mod.ProductFiles.Ctl.ProductFilesDelete)
		}
		categories := auth.Group("/categories")
		{
			categories.GET("/", mod.Categories.Ctl.CategoriesList)
			categories.GET("/:id", mod.Categories.Ctl.CategoriesInfo)
			categories.POST("/", mod.Categories.Ctl.CreateCategoryController)
			categories.PATCH("/:id", mod.Categories.Ctl.CategoriesUpdate)
			categories.DELETE("/:id", mod.Categories.Ctl.CategoriesDelete)
		}
	}
}

func apiProductPublic(r *gin.RouterGroup, mod *modules.Modules) {
	Public := r.Group("/public")
	{
		products := Public.Group("/products")
		{
			products.GET("/", mod.Products.Ctl.ProductsList)
			products.GET("/:id", mod.Products.Ctl.ProductsInfo)
		}
		product_details := Public.Group("/product_details")
		{
			product_details.GET("/", mod.ProductDetails.Ctl.ProductDetailsList)
			product_details.GET("/:id", mod.ProductDetails.Ctl.ProductDetailsInfo)
		}
		product_files := Public.Group("/product_files")
		{
			product_files.GET("/", mod.ProductFiles.Ctl.ProductFilesList)
			product_files.GET("/:id", mod.ProductFiles.Ctl.ProductFilesInfo)
		}
		categories := Public.Group("/categories")
		{
			categories.GET("/", mod.Categories.Ctl.CategoriesList)
			categories.GET("/:id", mod.Categories.Ctl.CategoriesInfo)
		}
	}
}

func apiOrder(r *gin.RouterGroup, mod *modules.Modules) {
	auth := r.Group("/auth", mod.Auth.Ctl.AuthMiddleware())
	{
		orders := auth.Group("/orders")
		{
			orders.GET("/", mod.Orders.Ctl.OrdersList)
			orders.GET("/:id", mod.Orders.Ctl.OrdersInfo)
			orders.POST("/", mod.Orders.Ctl.CreateOrderController)
			orders.PATCH("/:id", mod.Orders.Ctl.OrdersUpdate)
			orders.DELETE("/:id", mod.Orders.Ctl.OrdersDelete)
		}
		order_items := auth.Group("/order_items")
		{
			order_items.GET("/", mod.OrderItems.Ctl.OrderItemsList)
			order_items.GET("/:id", mod.OrderItems.Ctl.OrderItemsInfo)
			order_items.POST("/", mod.OrderItems.Ctl.CreateOrderItemController)
			order_items.PATCH("/:id", mod.OrderItems.Ctl.OrderItemsUpdate)
			order_items.DELETE("/:id", mod.OrderItems.Ctl.OrderItemsDelete)
		}
		carts := auth.Group("/carts")
		{
			carts.GET("/", mod.Carts.Ctl.CartsList)
			carts.GET("/:id", mod.Carts.Ctl.CartsInfo)
			carts.POST("/", mod.Carts.Ctl.CreateCartController)
			carts.PATCH("/:id", mod.Carts.Ctl.CartsUpdate)
			carts.DELETE("/:id", mod.Carts.Ctl.CartsDelete)
		}
		cart_items := auth.Group("/cart_items")
		{
			cart_items.GET("/", mod.CartItems.Ctl.CartItemsList)
			cart_items.GET("/:id", mod.CartItems.Ctl.CartItemsInfo)
			cart_items.POST("/", mod.CartItems.Ctl.CreateCartItemController)
			cart_items.POST("/selection", mod.CartItems.Ctl.SelectItemsController)
			cart_items.PATCH("/:id", mod.CartItems.Ctl.CartItemsUpdate)
			cart_items.DELETE("/:id", mod.CartItems.Ctl.CartItemsDelete)
		}
		payments := auth.Group("/payments")
		{
			payments.GET("/", mod.Payments.Ctl.PaymentsList)
			payments.GET("/:id", mod.Payments.Ctl.PaymentsInfo)
			payments.POST("/", mod.Payments.Ctl.CreatePaymentController)
			payments.PATCH("/:id", mod.Payments.Ctl.PaymentsUpdate)
			payments.DELETE("/:id", mod.Payments.Ctl.PaymentsDelete)
		}
		payment_files := auth.Group("/payment_files")
		{
			payment_files.GET("/", mod.PaymentFiles.Ctl.PaymentFilesList)
			payment_files.GET("/:id", mod.PaymentFiles.Ctl.PaymentFilesInfo)
			payment_files.POST("/", mod.PaymentFiles.Ctl.CreatePaymentFileController)
			payment_files.PATCH("/:id", mod.PaymentFiles.Ctl.PaymentFilesUpdate)
			payment_files.DELETE("/:id", mod.PaymentFiles.Ctl.PaymentFilesDelete)
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
