package members

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberAddressURIRequest struct {
	MemberID  string `uri:"id"`
	AddressID string `uri:"address_id"`
}

type CreateMemberAddressControllerRequest struct {
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

type UpdateMemberAddressControllerRequest = CreateMemberAddressControllerRequest

type ListMemberAddressesControllerRequest struct {
	base.RequestPaginate
}

func (c *Controller) ListMemberAddressesController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.address.list.start`)

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req ListMemberAddressesControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListMemberAddressesService(ctx.Request.Context(), &ListMemberAddressesServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.address.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) CreateMemberAddressController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.address.create.start`)

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req CreateMemberAddressControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	subDistrictID, err := uuid.Parse(req.SubDistrictID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	districtID, err := uuid.Parse(req.DistrictID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	provinceID, err := uuid.Parse(req.ProvinceID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	zipcodeID, err := uuid.Parse(req.ZipcodeID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.CreateMemberAddressService(ctx.Request.Context(), memberID, &CreateMemberAddressServiceRequest{
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
		ActionBy:      actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.address.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) InfoMemberAddressController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.address.info.start`)

	memberID, addressID, ok := c.parseMemberAddressURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	data, err := c.svc.InfoMemberAddressService(ctx.Request.Context(), memberID, addressID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.address.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) UpdateMemberAddressController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.address.update.start`)

	memberID, addressID, ok := c.parseMemberAddressURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req UpdateMemberAddressControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	subDistrictID, err := uuid.Parse(req.SubDistrictID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	districtID, err := uuid.Parse(req.DistrictID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	provinceID, err := uuid.Parse(req.ProvinceID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	zipcodeID, err := uuid.Parse(req.ZipcodeID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.UpdateMemberAddressService(ctx.Request.Context(), memberID, addressID, &UpdateMemberAddressServiceRequest{
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
		ActionBy:      actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.address.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteMemberAddressController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.address.delete.start`)

	memberID, addressID, ok := c.parseMemberAddressURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.DeleteMemberAddressService(ctx.Request.Context(), memberID, addressID, actionBy); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.address.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseMemberAddressURI(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	var uri MemberAddressURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	memberID, err := uuid.Parse(uri.MemberID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	addressID, err := uuid.Parse(uri.AddressID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return memberID, addressID, true
}

func (c *Controller) parseMemberID(ctx *gin.Context) (uuid.UUID, bool) {
	id := ctx.Param("id")
	memberID, err := uuid.Parse(id)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, false
	}
	return memberID, true
}

func getActionBy(ctx *gin.Context) *uuid.UUID {
	if memberID, ok := auth.GetMemberID(ctx); ok {
		return &memberID
	}
	return nil
}
