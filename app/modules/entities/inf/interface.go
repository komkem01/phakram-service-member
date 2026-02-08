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
