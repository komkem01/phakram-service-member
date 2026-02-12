package members

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterControllerRequest struct {
	PrefixID    string `json:"prefix_id"`
	GenderID    string `json:"gender_id"`
	FirstnameTh string `json:"firstname_th"`
	LastnameTh  string `json:"lastname_th"`
	FirstnameEn string `json:"firstname_en"`
	LastnameEn  string `json:"lastname_en"`
	Role        string `json:"role"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	TierID      string `json:"tier_id"`
	StatusID    string `json:"status_id"`
}

func (c *Controller) CreateRegisterController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`register.ctl.create.start`)

	var req RegisterControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`register.ctl.create.request`)

	prefixID, err := uuid.Parse(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var tierID uuid.UUID
	if req.TierID != "" {
		parsedTierID, err := uuid.Parse(req.TierID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		tierID = parsedTierID
	}

	var statusID uuid.UUID
	if req.StatusID != "" {
		parsedStatusID, err := uuid.Parse(req.StatusID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		statusID = parsedStatusID
	}

	if err := c.svc.CreateRegisterService(ctx.Request.Context(), &RegisterServiceRequest{
		PrefixID:    prefixID,
		GenderID:    genderID,
		FirstnameTh: req.FirstnameTh,
		LastnameTh:  req.LastnameTh,
		FirstnameEn: req.FirstnameEn,
		LastnameEn:  req.LastnameEn,
		Role:        req.Role,
		Phone:       req.Phone,
		Email:       req.Email,
		Password:    req.Password,
		TierID:      tierID,
		StatusID:    statusID,
		ActionBy:    nil,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`register.ctl.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) CreateRegisterByAdminController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`register.ctl.create_by_admin.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var req RegisterControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`register.ctl.create_by_admin.request`)

	prefixID, err := uuid.Parse(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var actionBy *uuid.UUID
	if memberID, ok := auth.GetMemberID(ctx); ok {
		actionBy = &memberID
	}

	var adminTierID uuid.UUID
	if req.TierID != "" {
		parsedTierID, err := uuid.Parse(req.TierID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		adminTierID = parsedTierID
	}

	var adminStatusID uuid.UUID
	if req.StatusID != "" {
		parsedStatusID, err := uuid.Parse(req.StatusID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		adminStatusID = parsedStatusID
	}

	if err := c.svc.CreateRegisterByAdminService(ctx.Request.Context(), &RegisterServiceRequest{
		PrefixID:    prefixID,
		GenderID:    genderID,
		FirstnameTh: req.FirstnameTh,
		LastnameTh:  req.LastnameTh,
		FirstnameEn: req.FirstnameEn,
		LastnameEn:  req.LastnameEn,
		Role:        req.Role,
		Phone:       req.Phone,
		Email:       req.Email,
		Password:    req.Password,
		TierID:      adminTierID,
		StatusID:    adminStatusID,
		ActionBy:    actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`register.ctl.create_by_admin.success`)
	base.Success(ctx, nil)
}
