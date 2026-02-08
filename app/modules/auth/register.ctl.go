package auth

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RegisterMemberController struct {
	MemberNo      string          `json:"member_no"`
	TierID        string          `json:"tier_id"`
	StatusID      string          `json:"status_id"`
	PrefixID      string          `json:"prefix_id"`
	GenderID      string          `json:"gender_id"`
	FirstnameTh   string          `json:"firstname_th"`
	LastnameTh    string          `json:"lastname_th"`
	FirstnameEn   string          `json:"firstname_en"`
	LastnameEn    string          `json:"lastname_en"`
	Role          string          `json:"role"`
	Phone         string          `json:"phone"`
	Email         string          `json:"email"`
	Password      string          `json:"password"`
	TotalSpent    decimal.Decimal `json:"total_spent"`
	CurrentPoints int             `json:"current_points"`
	Registration  time.Time       `json:"registration"`
	LastLogin     time.Time       `json:"last_login"`
}

func (c *Controller) RegisterMemberController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.register.start`)

	var req RegisterMemberController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`members.ctl.register.request`)

	tierID, err := uuid.Parse(req.TierID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	statusID, err := uuid.Parse(req.StatusID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	prefixID, err := uuid.Parse(req.PrefixID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.RegisterMemberService(ctx.Request.Context(), &RegisterMemberService{
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
		Email:         req.Email,
		Password:      req.Password,
		TotalSpent:    req.TotalSpent,
		CurrentPoints: req.CurrentPoints,
		Registration:  req.Registration,
		LastLogin:     req.LastLogin,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.register.success`)
	base.Success(ctx, nil)
}
