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

type UpdateMemberAddressControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateMemberAddressController struct {
	MemberID      string `json:"member_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	IsDefault     *bool  `json:"is_default"`
	AddressNo     string `json:"address_no"`
	Village       string `json:"village"`
	Alley         string `json:"alley"`
	SubDistrictID string `json:"sub_district_id"`
	DistrictID    string `json:"district_id"`
	ProvinceID    string `json:"province_id"`
	ZipcodeID     string `json:"zipcode_id"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateMemberAddressControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_addresses.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateMemberAddressController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_addresses.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}
	var subDistrictID uuid.UUID
	if req.SubDistrictID != "" {
		subDistrictID, err = uuid.Parse(req.SubDistrictID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var districtID uuid.UUID
	if req.DistrictID != "" {
		districtID, err = uuid.Parse(req.DistrictID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var provinceID uuid.UUID
	if req.ProvinceID != "" {
		provinceID, err = uuid.Parse(req.ProvinceID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var zipcodeID uuid.UUID
	if req.ZipcodeID != "" {
		zipcodeID, err = uuid.Parse(req.ZipcodeID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateMemberAddressService{
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
	span.AddEvent(`member_addresses.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) MemberAddressesUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
