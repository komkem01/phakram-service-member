package orders

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderItemURIRequest struct {
	OrderID string `uri:"id"`
	ItemID  string `uri:"item_id"`
}

type ListOrderItemControllerRequest struct {
	base.RequestPaginate
}

type CreateOrderItemControllerRequest struct {
	ProductID       string `json:"product_id"`
	Quantity        int    `json:"quantity"`
	PricePerUnit    string `json:"price_per_unit"`
	TotalItemAmount string `json:"total_item_amount"`
}

type UpdateOrderItemControllerRequest = CreateOrderItemControllerRequest

func (c *Controller) ListOrderItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.items.list.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	var req ListOrderItemControllerRequest
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

	data, page, err := c.svc.ListOrderItemService(ctx.Request.Context(), &ListOrderItemServiceRequest{RequestPaginate: req.RequestPaginate, OrderID: orderID}, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.items.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) InfoOrderItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.items.info.start`)

	orderID, itemID, ok := c.parseOrderItemID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.InfoOrderItemService(ctx.Request.Context(), orderID, itemID, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.items.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) CreateOrderItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.items.create.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	var req CreateOrderItemControllerRequest
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

	if err := c.svc.CreateOrderItemService(ctx.Request.Context(), orderID, &CreateOrderItemServiceRequest{
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
	}, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.items.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) UpdateOrderItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.items.update.start`)

	orderID, itemID, ok := c.parseOrderItemID(ctx)
	if !ok {
		return
	}

	var req UpdateOrderItemControllerRequest
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

	if err := c.svc.UpdateOrderItemService(ctx.Request.Context(), orderID, itemID, &UpdateOrderItemServiceRequest{
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
	}, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.items.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteOrderItemController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.items.delete.start`)

	orderID, itemID, ok := c.parseOrderItemID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	if err := c.svc.DeleteOrderItemService(ctx.Request.Context(), orderID, itemID, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.items.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseOrderItemID(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	var uri OrderItemURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	orderID, err := uuid.Parse(uri.OrderID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	itemID, err := uuid.Parse(uri.ItemID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return orderID, itemID, true
}
