package entitiesinf

import (
	"context"

	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

// ObjectEntity defines the interface for object entity operations such as create, retrieve, update, and soft delete.
type ExampleEntity interface {
	CreateExample(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
	GetExampleByID(ctx context.Context, id uuid.UUID) (*ent.Example, error)
	UpdateExampleByID(ctx context.Context, id uuid.UUID, status ent.ExampleStatus) (*ent.Example, error)
	SoftDeleteExampleByID(ctx context.Context, id uuid.UUID) error
	ListExamplesByStatus(ctx context.Context, status ent.ExampleStatus) ([]*ent.Example, error)
}
type ExampleTwoEntity interface {
	CreateExampleTwo(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
}

type GenderEntity interface {
	ListGenders(ctx context.Context, req *entitiesdto.ListGendersRequest) ([]*ent.GenderEntity, *base.ResponsePaginate, error)
	GetGenderByID(ctx context.Context, id uuid.UUID) (*ent.GenderEntity, error)
	CreateGender(ctx context.Context, gender *ent.GenderEntity) error
	UpdateGender(ctx context.Context, gender *ent.GenderEntity) error
	DeleteGender(ctx context.Context, genderID uuid.UUID) error
}

type PrefixEntity interface {
	ListPrefixes(ctx context.Context, req *entitiesdto.ListPrefixesRequest) ([]*ent.PrefixEntity, *base.ResponsePaginate, error)
	GetPrefixByID(ctx context.Context, id uuid.UUID) (*ent.PrefixEntity, error)
	CreatePrefix(ctx context.Context, prefix *ent.PrefixEntity) error
	UpdatePrefix(ctx context.Context, prefix *ent.PrefixEntity) error
	DeletePrefix(ctx context.Context, prefixID uuid.UUID) error
}

type ProvinceEntity interface {
	ListProvinces(ctx context.Context, req *entitiesdto.ListProvincesRequest) ([]*ent.ProvinceEntity, *base.ResponsePaginate, error)
	GetProvinceByID(ctx context.Context, id uuid.UUID) (*ent.ProvinceEntity, error)
	CreateProvince(ctx context.Context, province *ent.ProvinceEntity) error
	UpdateProvince(ctx context.Context, province *ent.ProvinceEntity) error
	DeleteProvince(ctx context.Context, provinceID uuid.UUID) error
}

type DistrictEntity interface {
	ListDistricts(ctx context.Context, req *entitiesdto.ListDistrictsRequest) ([]*ent.DistrictEntity, *base.ResponsePaginate, error)
	GetDistrictByID(ctx context.Context, id uuid.UUID) (*ent.DistrictEntity, error)
	CreateDistrict(ctx context.Context, district *ent.DistrictEntity) error
	UpdateDistrict(ctx context.Context, district *ent.DistrictEntity) error
	DeleteDistrict(ctx context.Context, districtID uuid.UUID) error
}

type SubDistrictEntity interface {
	ListSubDistricts(ctx context.Context, req *entitiesdto.ListSubDistrictsRequest) ([]*ent.SubDistrictEntity, *base.ResponsePaginate, error)
	GetSubDistrictByID(ctx context.Context, id uuid.UUID) (*ent.SubDistrictEntity, error)
	CreateSubDistrict(ctx context.Context, subDistrict *ent.SubDistrictEntity) error
	UpdateSubDistrict(ctx context.Context, subDistrict *ent.SubDistrictEntity) error
	DeleteSubDistrict(ctx context.Context, subDistrictID uuid.UUID) error
}

type ZipcodeEntity interface {
	ListZipcodes(ctx context.Context, req *entitiesdto.ListZipcodesRequest) ([]*ent.ZipcodeEntity, *base.ResponsePaginate, error)
	GetZipcodeByID(ctx context.Context, id uuid.UUID) (*ent.ZipcodeEntity, error)
	CreateZipcode(ctx context.Context, zipcode *ent.ZipcodeEntity) error
	UpdateZipcode(ctx context.Context, zipcode *ent.ZipcodeEntity) error
	DeleteZipcode(ctx context.Context, zipcodeID uuid.UUID) error
}

type StatusEntity interface {
	ListStatuses(ctx context.Context, req *entitiesdto.ListStatusesRequest) ([]*ent.StatusEntity, *base.ResponsePaginate, error)
	GetStatusByID(ctx context.Context, id uuid.UUID) (*ent.StatusEntity, error)
	CreateStatus(ctx context.Context, status *ent.StatusEntity) error
	UpdateStatus(ctx context.Context, status *ent.StatusEntity) error
	DeleteStatus(ctx context.Context, statusID uuid.UUID) error
}

type TierEntity interface {
	ListTiers(ctx context.Context, req *entitiesdto.ListTiersRequest) ([]*ent.TierEntity, *base.ResponsePaginate, error)
	GetTierByID(ctx context.Context, id uuid.UUID) (*ent.TierEntity, error)
	CreateTier(ctx context.Context, tier *ent.TierEntity) error
	UpdateTier(ctx context.Context, tier *ent.TierEntity) error
	DeleteTier(ctx context.Context, tierID uuid.UUID) error
}

type BankEntity interface {
	ListBanks(ctx context.Context, req *entitiesdto.ListBanksRequest) ([]*ent.BankEntity, *base.ResponsePaginate, error)
	GetBankByID(ctx context.Context, id uuid.UUID) (*ent.BankEntity, error)
	CreateBank(ctx context.Context, bank *ent.BankEntity) error
	UpdateBank(ctx context.Context, bank *ent.BankEntity) error
	DeleteBank(ctx context.Context, bankID uuid.UUID) error
}

type StorageEntity interface {
	ListStorages(ctx context.Context, req *entitiesdto.ListStoragesRequest) ([]*ent.StorageEntity, *base.ResponsePaginate, error)
	GetStorageByID(ctx context.Context, id uuid.UUID) (*ent.StorageEntity, error)
	CreateStorage(ctx context.Context, storage *ent.StorageEntity) error
	UpdateStorage(ctx context.Context, storage *ent.StorageEntity) error
	DeleteStorage(ctx context.Context, storageID uuid.UUID) error
}

type MemberEntity interface {
	ListMembers(ctx context.Context, req *entitiesdto.ListMembersRequest) ([]*ent.MemberEntity, *base.ResponsePaginate, error)
	GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error)
	GetMemberByPhone(ctx context.Context, phone string) (*ent.MemberEntity, error)
	CreateMember(ctx context.Context, member *ent.MemberEntity) error
	UpdateMember(ctx context.Context, member *ent.MemberEntity) error
	DeleteMember(ctx context.Context, memberID uuid.UUID) error
}

type MemberAccountEntity interface {
	ListMemberAccounts(ctx context.Context, req *entitiesdto.ListMemberAccountsRequest) ([]*ent.MemberAccountEntity, *base.ResponsePaginate, error)
	GetMemberAccountByID(ctx context.Context, id uuid.UUID) (*ent.MemberAccountEntity, error)
	GetMemberAccountByEmail(ctx context.Context, email string) (*ent.MemberAccountEntity, error)
	CreateMemberAccount(ctx context.Context, account *ent.MemberAccountEntity) error
	UpdateMemberAccount(ctx context.Context, account *ent.MemberAccountEntity) error
	DeleteMemberAccount(ctx context.Context, accountID uuid.UUID) error
	GetMemberAccountByMemberID(ctx context.Context, memberID uuid.UUID) (*ent.MemberAccountEntity, error)
}

type MemberAddressEntity interface {
	ListMemberAddresses(ctx context.Context, req *entitiesdto.ListMemberAddressesRequest) ([]*ent.MemberAddressEntity, *base.ResponsePaginate, error)
	GetMemberAddressByID(ctx context.Context, id uuid.UUID) (*ent.MemberAddressEntity, error)
	CreateMemberAddress(ctx context.Context, address *ent.MemberAddressEntity) error
	UpdateMemberAddress(ctx context.Context, address *ent.MemberAddressEntity) error
	DeleteMemberAddress(ctx context.Context, addressID uuid.UUID) error
}

type MemberBankEntity interface {
	ListMemberBanks(ctx context.Context, req *entitiesdto.ListMemberBanksRequest) ([]*ent.MemberBankEntity, *base.ResponsePaginate, error)
	GetMemberBankByID(ctx context.Context, id uuid.UUID) (*ent.MemberBankEntity, error)
	CreateMemberBank(ctx context.Context, bank *ent.MemberBankEntity) error
	UpdateMemberBank(ctx context.Context, bank *ent.MemberBankEntity) error
	DeleteMemberBank(ctx context.Context, bankID uuid.UUID) error
}

type MemberFileEntity interface {
	ListMemberFiles(ctx context.Context, req *entitiesdto.ListMemberFilesRequest) ([]*ent.MemberFileEntity, *base.ResponsePaginate, error)
	GetMemberFileByID(ctx context.Context, id uuid.UUID) (*ent.MemberFileEntity, error)
	GetMemberFileByMemberID(ctx context.Context, memberID uuid.UUID) (*ent.MemberFileEntity, error)
	CreateMemberFile(ctx context.Context, file *ent.MemberFileEntity) error
	UpdateMemberFile(ctx context.Context, file *ent.MemberFileEntity) error
	DeleteMemberFile(ctx context.Context, fileID uuid.UUID) error
}

type MemberTransactionEntity interface {
	ListMemberTransactions(ctx context.Context, req *entitiesdto.ListMemberTransactionsRequest) ([]*ent.MemberTransactionEntity, *base.ResponsePaginate, error)
	GetMemberTransactionByID(ctx context.Context, id uuid.UUID) (*ent.MemberTransactionEntity, error)
	CreateMemberTransaction(ctx context.Context, transaction *ent.MemberTransactionEntity) error
	UpdateMemberTransaction(ctx context.Context, transaction *ent.MemberTransactionEntity) error
	DeleteMemberTransaction(ctx context.Context, transactionID uuid.UUID) error
}

type AuditLogEntity interface {
	ListAuditLogs(ctx context.Context, req *entitiesdto.ListAuditLogsRequest) ([]*ent.AuditLogEntity, *base.ResponsePaginate, error)
	GetAuditLogByID(ctx context.Context, id uuid.UUID) (*ent.AuditLogEntity, error)
	CreateAuditLog(ctx context.Context, log *ent.AuditLogEntity) error
	UpdateAuditLog(ctx context.Context, log *ent.AuditLogEntity) error
	DeleteAuditLog(ctx context.Context, logID uuid.UUID) error
}

type MemberWishlistEntity interface {
	ListMemberWishlist(ctx context.Context, req *entitiesdto.ListMemberWishlistRequest) ([]*ent.MemberWishlistEntity, *base.ResponsePaginate, error)
	GetMemberWishlistByID(ctx context.Context, id uuid.UUID) (*ent.MemberWishlistEntity, error)
	CreateMemberWishlist(ctx context.Context, wishlist *ent.MemberWishlistEntity) error
	UpdateMemberWishlist(ctx context.Context, wishlist *ent.MemberWishlistEntity) error
	DeleteMemberWishlist(ctx context.Context, wishlistID uuid.UUID) error
}

type ProductEntity interface {
	ListProducts(ctx context.Context, req *entitiesdto.ListProductsRequest) ([]*ent.ProductEntity, *base.ResponsePaginate, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*ent.ProductEntity, error)
	CreateProduct(ctx context.Context, product *ent.ProductEntity) error
	UpdateProduct(ctx context.Context, product *ent.ProductEntity) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
}

type ProductDetailEntity interface {
	ListProductDetails(ctx context.Context, req *entitiesdto.ListProductDetailsRequest) ([]*ent.ProductDetailEntity, *base.ResponsePaginate, error)
	GetProductDetailByID(ctx context.Context, id uuid.UUID) (*ent.ProductDetailEntity, error)
	CreateProductDetail(ctx context.Context, detail *ent.ProductDetailEntity) error
	UpdateProductDetail(ctx context.Context, detail *ent.ProductDetailEntity) error
	DeleteProductDetail(ctx context.Context, detailID uuid.UUID) error
}

type ProductStockEntity interface {
	ListProductStocks(ctx context.Context, req *entitiesdto.ListProductStocksRequest) ([]*ent.ProductStockEntity, *base.ResponsePaginate, error)
	GetProductStockByID(ctx context.Context, id uuid.UUID) (*ent.ProductStockEntity, error)
	CreateProductStock(ctx context.Context, stock *ent.ProductStockEntity) error
	UpdateProductStock(ctx context.Context, stock *ent.ProductStockEntity) error
	DeleteProductStock(ctx context.Context, stockID uuid.UUID) error
}

type ProductFileEntity interface {
	ListProductFiles(ctx context.Context, req *entitiesdto.ListProductFilesRequest) ([]*ent.ProductFileEntity, *base.ResponsePaginate, error)
	GetProductFileByID(ctx context.Context, id uuid.UUID) (*ent.ProductFileEntity, error)
	GetProductFileByProductID(ctx context.Context, productID uuid.UUID) (*ent.ProductFileEntity, error)
	CreateProductFile(ctx context.Context, file *ent.ProductFileEntity) error
	UpdateProductFile(ctx context.Context, file *ent.ProductFileEntity) error
	DeleteProductFile(ctx context.Context, fileID uuid.UUID) error
}

type CategoryEntity interface {
	ListCategories(ctx context.Context, req *entitiesdto.ListCategoriesRequest) ([]*ent.CategoryEntity, *base.ResponsePaginate, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*ent.CategoryEntity, error)
	CreateCategory(ctx context.Context, category *ent.CategoryEntity) error
	UpdateCategory(ctx context.Context, category *ent.CategoryEntity) error
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
}

type OrderEntity interface {
	ListOrders(ctx context.Context, req *entitiesdto.ListOrdersRequest) ([]*ent.OrderEntity, *base.ResponsePaginate, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*ent.OrderEntity, error)
	CreateOrder(ctx context.Context, order *ent.OrderEntity) error
	UpdateOrder(ctx context.Context, order *ent.OrderEntity) error
	DeleteOrder(ctx context.Context, orderID uuid.UUID) error
}

type OrderItemEntity interface {
	ListOrderItems(ctx context.Context, req *entitiesdto.ListOrderItemsRequest) ([]*ent.OrderItemEntity, *base.ResponsePaginate, error)
	GetOrderItemByID(ctx context.Context, id uuid.UUID) (*ent.OrderItemEntity, error)
	CreateOrderItem(ctx context.Context, item *ent.OrderItemEntity) error
	UpdateOrderItem(ctx context.Context, item *ent.OrderItemEntity) error
	DeleteOrderItem(ctx context.Context, itemID uuid.UUID) error
}

type CartEntity interface {
	ListCarts(ctx context.Context, req *entitiesdto.ListCartsRequest) ([]*ent.CartEntity, *base.ResponsePaginate, error)
	GetCartByID(ctx context.Context, id uuid.UUID) (*ent.CartEntity, error)
	CreateCart(ctx context.Context, cart *ent.CartEntity) error
	UpdateCart(ctx context.Context, cart *ent.CartEntity) error
	DeleteCart(ctx context.Context, cartID uuid.UUID) error
}

type CartItemEntity interface {
	ListCartItems(ctx context.Context, req *entitiesdto.ListCartItemsRequest) ([]*ent.CartItemEntity, *base.ResponsePaginate, error)
	GetCartItemByID(ctx context.Context, id uuid.UUID) (*ent.CartItemEntity, error)
	CreateCartItem(ctx context.Context, item *ent.CartItemEntity) error
	UpdateCartItem(ctx context.Context, item *ent.CartItemEntity) error
	DeleteCartItem(ctx context.Context, itemID uuid.UUID) error
}

type PaymentEntity interface {
	ListPayments(ctx context.Context, req *entitiesdto.ListPaymentsRequest) ([]*ent.PaymentEntity, *base.ResponsePaginate, error)
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*ent.PaymentEntity, error)
	CreatePayment(ctx context.Context, payment *ent.PaymentEntity) error
	UpdatePayment(ctx context.Context, payment *ent.PaymentEntity) error
	DeletePayment(ctx context.Context, paymentID uuid.UUID) error
}

type PaymentFileEntity interface {
	ListPaymentFiles(ctx context.Context, req *entitiesdto.ListPaymentFilesRequest) ([]*ent.PaymentFileEntity, *base.ResponsePaginate, error)
	GetPaymentFileByID(ctx context.Context, id uuid.UUID) (*ent.PaymentFileEntity, error)
	CreatePaymentFile(ctx context.Context, file *ent.PaymentFileEntity) error
	UpdatePaymentFile(ctx context.Context, file *ent.PaymentFileEntity) error
	DeletePaymentFile(ctx context.Context, fileID uuid.UUID) error
}
