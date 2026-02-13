package carts

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListCartControllerRequest struct {
	base.RequestPaginate
	MemberID string `form:"member_id"`
}

type CreateCartControllerRequest struct {
	MemberID string `json:"member_id"`
	IsActive *bool  `json:"is_active"`
}

type UpdateCartControllerRequest struct {
	IsActive *bool `json:"is_active"`
}

type CartURIRequest struct {
	ID string `uri:"id"`
}

func (c *Controller) ListCartController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.list.start`)

	var req ListCartControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID := uuid.Nil
	if isAdmin {
		if req.MemberID != "" {
			parsedMemberID, err := uuid.Parse(req.MemberID)
			if err != nil {
				base.BadRequest(ctx, i18n.BadRequest, nil)
				return
			}
			memberID = parsedMemberID
		}
	} else {
		memberID = requesterID
	}

	data, page, err := c.svc.ListCartService(ctx.Request.Context(), &ListCartServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) InfoCartController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.info.start`)

	cartID, ok := c.parseCartID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.InfoCartService(ctx.Request.Context(), cartID, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) CreateCartController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.create.start`)

	var req CreateCartControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID := requesterID
	if isAdmin {
		if req.MemberID == "" {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		parsedMemberID, err := uuid.Parse(req.MemberID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		memberID = parsedMemberID
	}

	if err := c.svc.CreateCartService(ctx.Request.Context(), &CreateCartServiceRequest{MemberID: memberID, IsActive: req.IsActive}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) UpdateCartController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.update.start`)

	cartID, ok := c.parseCartID(ctx)
	if !ok {
		return
	}

	var req UpdateCartControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	if err := c.svc.UpdateCartService(ctx.Request.Context(), cartID, &UpdateCartServiceRequest{IsActive: req.IsActive}, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteCartController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.delete.start`)

	cartID, ok := c.parseCartID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	if err := c.svc.DeleteCartService(ctx.Request.Context(), cartID, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseCartID(ctx *gin.Context) (uuid.UUID, bool) {
	var uri CartURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, false
	}

	cartID, err := uuid.Parse(uri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, false
	}

	return cartID, true
}
