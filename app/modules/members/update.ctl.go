package members

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateMemberControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateMemberController struct {
	MemberNo      string           `json:"member_no"`
	TierID        string           `json:"tier_id"`
	StatusID      string           `json:"status_id"`
	PrefixID      string           `json:"prefix_id"`
	GenderID      string           `json:"gender_id"`
	FirstnameTh   string           `json:"firstname_th"`
	LastnameTh    string           `json:"lastname_th"`
	FirstnameEn   string           `json:"firstname_en"`
	LastnameEn    string           `json:"lastname_en"`
	Role          string           `json:"role"`
	Phone         string           `json:"phone"`
	TotalSpent    *decimal.Decimal `json:"total_spent"`
	CurrentPoints *int             `json:"current_points"`
	Registration  *time.Time       `json:"registration"`
	LastLogin     *time.Time       `json:"last_login"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateMemberControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`members.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateMemberController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`members.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var tierID uuid.UUID
	if req.TierID != "" {
		tierID, err = uuid.Parse(req.TierID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var statusID uuid.UUID
	if req.StatusID != "" {
		statusID, err = uuid.Parse(req.StatusID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var prefixID uuid.UUID
	if req.PrefixID != "" {
		prefixID, err = uuid.Parse(req.PrefixID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var genderID uuid.UUID
	if req.GenderID != "" {
		genderID, err = uuid.Parse(req.GenderID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateMemberService{
		MemberNo:      req.MemberNo,
		TierID:        tierID,
		StatusID:      statusID,
		PrefixID:      prefixID,
		GenderID:      genderID,
		FirstnameTh:   req.FirstnameTh,
		LastnameTh:    req.LastnameTh,
		FirstnameEn:   req.FirstnameEn,
		LastnameEn:    req.LastnameEn,
		Role:          req.Role,
		Phone:         req.Phone,
		TotalSpent:    req.TotalSpent,
		CurrentPoints: req.CurrentPoints,
		Registration:  req.Registration,
		LastLogin:     req.LastLogin,
		MemberID:      memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`members.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) MembersUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
