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

type InfoMemberAddressControllerRequestUri struct {
	ID string `uri:"id"`
}

type InfoMemberAddressControllerResponses struct {
	ID            uuid.UUID `json:"id"`
	MemberID      uuid.UUID `json:"member_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone"`
	IsDefault     bool      `json:"is_default"`
	AddressNo     string    `json:"address_no"`
	Village       string    `json:"village"`
	Alley         string    `json:"alley"`
	SubDistrictID uuid.UUID `json:"sub_district_id"`
	DistrictID    uuid.UUID `json:"district_id"`
	ProvinceID    uuid.UUID `json:"province_id"`
	ZipcodeID     uuid.UUID `json:"zipcode_id"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req InfoMemberAddressControllerRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_addresses.ctl.info.request`)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	memberID := uuid.Nil
	isAdmin := authmod.GetIsAdmin(ctx)
	if !isAdmin {
		var ok bool
		memberID, ok = authmod.GetMemberID(ctx)
		if !ok {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			return
		}
	}

	data, err := c.svc.InfoService(ctx, id, memberID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_addresses.ctl.info.callsvc`)

	var resp InfoMemberAddressControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, resp)
}

func (c *Controller) MemberAddressesInfo(ctx *gin.Context) {
	c.InfoController(ctx)
}
