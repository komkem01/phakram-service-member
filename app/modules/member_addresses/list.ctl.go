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

type ListMemberAddressControllerRequest struct {
	base.RequestPaginate
}

type ListMemberAddressControllerResponses struct {
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

func (c *Controller) MemberAddressesList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListMemberAddressControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_addresses.ctl.list.request`)

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

	data, page, err := c.svc.ListService(ctx, &ListMemberAddressServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`member_addresses.ctl.list.callsvc`)

	var resp []*ListMemberAddressControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
