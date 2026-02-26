package auth

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"
	"strings"

	"github.com/gin-gonic/gin"
)

type CookiePolicyControllerQuery struct {
	VisitorKey string `form:"visitor_key"`
}

type AcceptCookieConsentControllerRequest struct {
	VisitorKey string `json:"visitor_key"`
}

type CreateCookiePolicyVersionControllerRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *Controller) GetCookiePolicyPublicController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.cookie_policy.public_get.start")

	var query CookiePolicyControllerQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.GetCookiePolicyInfoService(ctx.Request.Context(), nil, query.VisitorKey)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.cookie_policy.public_get.success")
	base.Success(ctx, data)
}

func (c *Controller) AcceptCookiePolicyPublicController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.cookie_policy.public_accept.start")

	var req AcceptCookieConsentControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.AcceptCookieConsentService(ctx.Request.Context(), &AcceptCookieConsentServiceRequest{
		VisitorKey: req.VisitorKey,
		UserAgent:  strings.TrimSpace(ctx.GetHeader("User-Agent")),
	}, nil)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.cookie_policy.public_accept.success")
	base.Success(ctx, data)
}

func (c *Controller) GetCookiePolicyAuthController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.cookie_policy.auth_get.start")

	memberID, ok := GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var query CookiePolicyControllerQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.GetCookiePolicyInfoService(ctx.Request.Context(), &memberID, query.VisitorKey)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.cookie_policy.auth_get.success")
	base.Success(ctx, data)
}

func (c *Controller) AcceptCookiePolicyAuthController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.cookie_policy.auth_accept.start")

	memberID, ok := GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var req AcceptCookieConsentControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.AcceptCookieConsentService(ctx.Request.Context(), &AcceptCookieConsentServiceRequest{
		VisitorKey: req.VisitorKey,
		UserAgent:  strings.TrimSpace(ctx.GetHeader("User-Agent")),
	}, &memberID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.cookie_policy.auth_accept.success")
	base.Success(ctx, data)
}

func (c *Controller) ListCookiePolicyVersionsController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.cookie_policy.versions_list.start")

	if !GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.ListCookiePolicyVersionsService(ctx.Request.Context())
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.cookie_policy.versions_list.success")
	base.Success(ctx, data)
}

func (c *Controller) CreateCookiePolicyVersionController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("auth.ctl.cookie_policy.versions_create.start")

	if !GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, ok := GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var req CreateCookiePolicyVersionControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.CreateCookiePolicyVersionService(ctx.Request.Context(), &CreateCookiePolicyVersionServiceRequest{
		Title:   req.Title,
		Content: req.Content,
	}, memberID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("auth.ctl.cookie_policy.versions_create.success")
	base.Success(ctx, data)
}
