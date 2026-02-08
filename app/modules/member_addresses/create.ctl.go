package member_addresses

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateMemberAddressController struct {
	MemberID      string `json:"member_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	IsDefault     bool   `json:"is_default"`
	AddressNo     string `json:"address_no"`
	Village       string `json:"village"`
	Alley         string `json:"alley"`
	SubDistrictID string `json:"sub_district_id"`
	DistrictID    string `json:"district_id"`
	ProvinceID    string `json:"province_id"`
	ZipcodeID     string `json:"zipcode_id"`
}

func (c *Controller) CreateMemberAddressController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`member_addresses.ctl.create.start`)

	var req CreateMemberAddressController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_addresses.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}
	subDistrictID, err := uuid.Parse(req.SubDistrictID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	districtID, err := uuid.Parse(req.DistrictID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	provinceID, err := uuid.Parse(req.ProvinceID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	zipcodeID, err := uuid.Parse(req.ZipcodeID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateMemberAddressService(ctx.Request.Context(), &CreateMemberAddressService{
		MemberID:      memberID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Phone:         req.Phone,
		IsDefault:     req.IsDefault,
		AddressNo:     req.AddressNo,
		Village:       req.Village,
		Alley:         req.Alley,
		SubDistrictID: subDistrictID,
		DistrictID:    districtID,
		ProvinceID:    provinceID,
		ZipcodeID:     zipcodeID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`member_addresses.ctl.create.success`)
	base.Success(ctx, nil)
}
