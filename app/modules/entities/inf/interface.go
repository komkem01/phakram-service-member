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
	UploadStorage(ctx context.Context, storage *ent.StorageEntity) error
	GetStorageByID(ctx context.Context, id uuid.UUID) (*ent.StorageEntity, error)
	ListStoragesByRefID(ctx context.Context, refID uuid.UUID) ([]*ent.StorageEntity, error)
	DeleteStorageByID(ctx context.Context, id uuid.UUID) error
	DeleteStoragesByRefID(ctx context.Context, refID uuid.UUID) error
	UpdateStatusStorage(ctx context.Context, id uuid.UUID, req *ent.StorageEntity) error
}

type MemberEntity interface {
	ListMembers(ctx context.Context, req *entitiesdto.ListMembersRequest) ([]*ent.MemberEntity, *base.ResponsePaginate, error)
	GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error)
	CreateMember(ctx context.Context, member *ent.MemberEntity) error
	UpdateMember(ctx context.Context, member *ent.MemberEntity) error
	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	// admin service
	CreateAdminMember(ctx context.Context, member *ent.MemberEntity) error
	UpdateAdminMember(ctx context.Context, member *ent.MemberEntity) error
	DeleteAdminMember(ctx context.Context, memberID uuid.UUID) error
	GetAdminMemberByID(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error)

	CreateMemberByAdmin(ctx context.Context, member *ent.MemberEntity) error
	UpdateMemberByAdmin(ctx context.Context, member *ent.MemberEntity) error
	DeleteMemberByAdmin(ctx context.Context, memberID uuid.UUID) error
	GetMemberByIDByAdmin(ctx context.Context, id uuid.UUID) (*ent.MemberEntity, error)
}

type MemberTransactionEntity interface {
	CreateMemberTransaction(ctx context.Context, memberTransaction *ent.MemberTransactionEntity) error
}

type MemberBankEntity interface {
	ListMemberBanks(ctx context.Context, req *entitiesdto.ListMemberBanksRequest) ([]*ent.MemberBankEntity, *base.ResponsePaginate, error)
	CreateMemberBank(ctx context.Context, memberBank *ent.MemberBankEntity) error
	GetMemberBankByID(ctx context.Context, id uuid.UUID) (*ent.MemberBankEntity, error)
	UpdateMemberBank(ctx context.Context, memberBank *ent.MemberBankEntity) error
	DeleteMemberBank(ctx context.Context, memberBankID uuid.UUID) error
}

type MemberAddressEntity interface {
	ListMemberAddresses(ctx context.Context, req *entitiesdto.ListMemberAddressesRequest) ([]*ent.MemberAddressEntity, *base.ResponsePaginate, error)
	CreateMemberAddress(ctx context.Context, memberAddress *ent.MemberAddressEntity) error
	GetMemberAddressByID(ctx context.Context, id uuid.UUID) (*ent.MemberAddressEntity, error)
	UpdateMemberAddress(ctx context.Context, memberAddress *ent.MemberAddressEntity) error
	DeleteMemberAddress(ctx context.Context, memberAddressID uuid.UUID) error
}

type MemberAccountEntity interface {
	ListMemberAccounts(ctx context.Context, req *entitiesdto.ListMemberAccountsRequest) ([]*ent.MemberAccountEntity, *base.ResponsePaginate, error)
	CreateMemberAccount(ctx context.Context, memberAccount *ent.MemberAccountEntity) error
	GetMemberAccountByID(ctx context.Context, id uuid.UUID) (*ent.MemberAccountEntity, error)
	UpdateMemberAccount(ctx context.Context, memberAccount *ent.MemberAccountEntity) error
	DeleteMemberAccount(ctx context.Context, memberAccountID uuid.UUID) error
}

type MemberWishlistEntity interface {
	ListMemberWishlist(ctx context.Context, req *entitiesdto.ListMemberWishlistRequest) ([]*ent.MemberWishlistEntity, *base.ResponsePaginate, error)
	CreateMemberWishlist(ctx context.Context, memberWishlist *ent.MemberWishlistEntity) error
	GetMemberWishlistByID(ctx context.Context, id uuid.UUID) (*ent.MemberWishlistEntity, error)
	UpdateMemberWishlist(ctx context.Context, memberWishlist *ent.MemberWishlistEntity) error
	DeleteMemberWishlist(ctx context.Context, memberWishlistID uuid.UUID) error
}

type MemberFileEntity interface {
	ListMemberFiles(ctx context.Context, req *entitiesdto.ListMemberFilesRequest) ([]*ent.MemberFileEntity, *base.ResponsePaginate, error)
	CreateMemberFile(ctx context.Context, memberFile *ent.MemberFileEntity) error
	GetMemberFileByID(ctx context.Context, id uuid.UUID) (*ent.MemberFileEntity, error)
	UpdateMemberFile(ctx context.Context, memberFile *ent.MemberFileEntity) error
	DeleteMemberFile(ctx context.Context, memberFileID uuid.UUID) error
}

type PaymentEntity interface {
	ListPayments(ctx context.Context, req *entitiesdto.ListPaymentsRequest) ([]*ent.PaymentEntity, *base.ResponsePaginate, error)
	CreatePayment(ctx context.Context, payment *ent.PaymentEntity) error
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*ent.PaymentEntity, error)
	UpdatePayment(ctx context.Context, payment *ent.PaymentEntity) error
	DeletePayment(ctx context.Context, paymentID uuid.UUID) error
}

type MemberPaymentEntity interface {
	ListMemberPayments(ctx context.Context, req *entitiesdto.ListMemberPaymentsRequest) ([]*ent.MemberPaymentEntity, *base.ResponsePaginate, error)
	CreateMemberPayment(ctx context.Context, memberPayment *ent.MemberPaymentEntity) error
	GetMemberPaymentByID(ctx context.Context, id uuid.UUID) (*ent.MemberPaymentEntity, error)
	UpdateMemberPayment(ctx context.Context, memberPayment *ent.MemberPaymentEntity) error
	DeleteMemberPayment(ctx context.Context, memberPaymentID uuid.UUID) error
}

type CategoryEntity interface {
	ListCategories(ctx context.Context, req *entitiesdto.ListCategoriesRequest) ([]*ent.CategoryEntity, *base.ResponsePaginate, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*ent.CategoryEntity, error)
	CreateCategory(ctx context.Context, category *ent.CategoryEntity) error
	UpdateCategory(ctx context.Context, category *ent.CategoryEntity) error
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
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
	CreateProductFile(ctx context.Context, file *ent.ProductFileEntity) error
	UpdateProductFile(ctx context.Context, file *ent.ProductFileEntity) error
	DeleteProductFile(ctx context.Context, fileID uuid.UUID) error
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

type PaymentFileEntity interface {
	CreatePaymentFile(ctx context.Context, paymentFile *ent.PaymentFileEntity) error
	GetPaymentFileByID(ctx context.Context, id uuid.UUID) (*ent.PaymentFileEntity, error)
	UpdatePaymentFile(ctx context.Context, paymentFile *ent.PaymentFileEntity) error
	DeletePaymentFile(ctx context.Context, paymentFileID uuid.UUID) error
}

type AuditLogEntity interface {
	ListAuditLogs(ctx context.Context, req *entitiesdto.ListAuditLogsRequest) ([]*ent.AuditLogEntity, *base.ResponsePaginate, error)
	GetAuditLogByID(ctx context.Context, id uuid.UUID) (*ent.AuditLogEntity, error)
	CreateAuditLog(ctx context.Context, log *ent.AuditLogEntity) error
	UpdateAuditLog(ctx context.Context, log *ent.AuditLogEntity) error
	DeleteAuditLog(ctx context.Context, logID uuid.UUID) error
}
