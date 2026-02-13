package auth

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginControllerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenControllerRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ActAsMemberControllerRequest struct {
	MemberID string `json:"member_id"`
}

func (c *Controller) LoginController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.login.start")

	var req LoginControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	res, err := c.svc.LoginService(ctx.Request.Context(), &LoginServiceRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.login.success")
	base.Success(ctx, res)
}

func (c *Controller) RefreshTokenController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.refresh.start")

	var req RefreshTokenControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	res, err := c.svc.RefreshTokenService(ctx.Request.Context(), &RefreshTokenServiceRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.refresh.success")
	base.Success(ctx, res)
}

func (c *Controller) GetInfoController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.me.start")

	memberID, ok := GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	res, err := c.svc.GetInfoService(ctx.Request.Context(), memberID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.me.success")
	base.Success(ctx, res)
}

func (c *Controller) ActAsMemberController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.act_as.start")

	if !GetActorIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	actorID, ok := GetActorMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var req ActAsMemberControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	targetID, err := uuid.Parse(req.MemberID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	res, err := c.svc.ActAsMemberService(ctx.Request.Context(), &ActAsMemberServiceRequest{ActorMemberID: actorID, TargetMemberID: targetID})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.act_as.success")
	base.Success(ctx, res)
}

func (c *Controller) ExitActAsController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.exit_act_as.start")

	if !GetActorIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	actorID, ok := GetActorMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	res, err := c.svc.ExitActAsService(ctx.Request.Context(), actorID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.exit_act_as.success")
	base.Success(ctx, res)
}
