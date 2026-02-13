package carts

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartItemURIRequest struct {
	CartID string `uri:"id"`
	ItemID string `uri:"item_id"`
}

type ListCartItemControllerRequest struct {
	base.RequestPaginate
}

type CreateCartItemControllerRequest struct {
	ProductID       string `json:"product_id"`
	Quantity        int    `json:"quantity"`
	PricePerUnit    string `json:"price_per_unit"`
	TotalItemAmount string `json:"total_item_amount"`
}

type UpdateCartItemControllerRequest = CreateCartItemControllerRequest

func (c *Controller) ListCartItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.items.list.start`)

	cartID, ok := c.parseCartID(ctx)
	if !ok {
		return
	}

	var req ListCartItemControllerRequest
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

	data, page, err := c.svc.ListCartItemService(ctx.Request.Context(), &ListCartItemServiceRequest{RequestPaginate: req.RequestPaginate, CartID: cartID}, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.items.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) InfoCartItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.items.info.start`)

	cartID, itemID, ok := c.parseCartItemID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.InfoCartItemService(ctx.Request.Context(), cartID, itemID, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.items.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) CreateCartItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.items.create.start`)

	cartID, ok := c.parseCartID(ctx)
	if !ok {
		return
	}

	var req CreateCartItemControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	if err := c.svc.CreateCartItemService(ctx.Request.Context(), cartID, &CreateCartItemServiceRequest{
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
	}, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.items.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) UpdateCartItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.items.update.start`)

	cartID, itemID, ok := c.parseCartItemID(ctx)
	if !ok {
		return
	}

	var req UpdateCartItemControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	if err := c.svc.UpdateCartItemService(ctx.Request.Context(), cartID, itemID, &UpdateCartItemServiceRequest{
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
	}, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.items.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteCartItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`carts.ctl.items.delete.start`)

	cartID, itemID, ok := c.parseCartItemID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	if err := c.svc.DeleteCartItemService(ctx.Request.Context(), cartID, itemID, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`carts.ctl.items.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseCartItemID(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	var uri CartItemURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	cartID, err := uuid.Parse(uri.CartID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	itemID, err := uuid.Parse(uri.ItemID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return cartID, itemID, true
}
