package members

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateControllerRequestURI struct {
	ID string `uri:"id"`
}

type UpdateControllerRequest struct {
	TierID      *string `json:"tier_id"`
	StatusID    *string `json:"status_id"`
	PrefixID    *string `json:"prefix_id"`
	GenderID    *string `json:"gender_id"`
	FirstnameTh *string `json:"firstname_th"`
	LastnameTh  *string `json:"lastname_th"`
	FirstnameEn *string `json:"firstname_en"`
	LastnameEn  *string `json:"lastname_en"`
	Role        *string `json:"role"`
	Phone       *string `json:"phone"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.update.start`)

	var uri UpdateControllerRequestURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if !c.ensureAdminOrSelf(ctx, id) {
		return
	}

	var req UpdateControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var tierID *uuid.UUID
	if req.TierID != nil && *req.TierID != "" {
		parsed, err := uuid.Parse(*req.TierID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		tierID = &parsed
	}
	var statusID *uuid.UUID
	if req.StatusID != nil && *req.StatusID != "" {
		parsed, err := uuid.Parse(*req.StatusID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		statusID = &parsed
	}
	var prefixID *uuid.UUID
	if req.PrefixID != nil && *req.PrefixID != "" {
		parsed, err := uuid.Parse(*req.PrefixID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		prefixID = &parsed
	}
	var genderID *uuid.UUID
	if req.GenderID != nil && *req.GenderID != "" {
		parsed, err := uuid.Parse(*req.GenderID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		genderID = &parsed
	}

	var actionBy *uuid.UUID
	if memberID, ok := auth.GetMemberID(ctx); ok {
		actionBy = &memberID
	}

	if err := c.svc.UpdateService(ctx.Request.Context(), id, &UpdateServiceRequest{
		TierID:      tierID,
		StatusID:    statusID,
		PrefixID:    prefixID,
		GenderID:    genderID,
		FirstnameTh: req.FirstnameTh,
		LastnameTh:  req.LastnameTh,
		FirstnameEn: req.FirstnameEn,
		LastnameEn:  req.LastnameEn,
		Role:        req.Role,
		Phone:       req.Phone,
		ActionBy:    actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.update.success`)
	base.Success(ctx, nil)
}
