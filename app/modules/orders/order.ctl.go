package orders

import (
	"phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListOrderControllerRequest struct {
	base.RequestPaginate
	MemberID  string `form:"member_id"`
	Search    string `form:"search"`
	Status    string `form:"status"`
	StartDate int64  `form:"start_date"`
	EndDate   int64  `form:"end_date"`
}

type CreateOrderControllerRequest struct {
	MemberID           string `json:"member_id"`
	PaymentID          string `json:"payment_id"`
	AddressID          string `json:"address_id"`
	Status             string `json:"status"`
	ShippingTrackingNo string `json:"shipping_tracking_no"`
	TotalAmount        string `json:"total_amount"`
	DiscountAmount     string `json:"discount_amount"`
	NetAmount          string `json:"net_amount"`
}

type UpdateOrderControllerRequest = CreateOrderControllerRequest

type OrderURIRequest struct {
	ID string `uri:"id"`
}

func (c *Controller) ListOrderController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.list.start`)

	var req ListOrderControllerRequest
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

	data, page, err := c.svc.ListOrderService(ctx.Request.Context(), &ListOrderServiceRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
		Search:          req.Search,
		Status:          req.Status,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) InfoOrderController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.info.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.InfoOrderService(ctx.Request.Context(), orderID, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) TimelineOrderController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.timeline.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	data, err := c.svc.TimelineOrderService(ctx.Request.Context(), orderID, requesterID, isAdmin)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.timeline.success`)
	base.Success(ctx, data)
}

func (c *Controller) CreateOrderController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.create.start`)

	var req CreateOrderControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	paymentID := uuid.Nil
	if req.PaymentID != "" {
		parsedPaymentID, err := uuid.Parse(req.PaymentID)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		paymentID = parsedPaymentID
	}
	addressID, err := uuid.Parse(req.AddressID)
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

	data, err := c.svc.CreateOrderService(ctx.Request.Context(), &CreateOrderServiceRequest{
		MemberID:           memberID,
		PaymentID:          paymentID,
		AddressID:          addressID,
		Status:             req.Status,
		ShippingTrackingNo: req.ShippingTrackingNo,
		TotalAmount:        req.TotalAmount,
		DiscountAmount:     req.DiscountAmount,
		NetAmount:          req.NetAmount,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.create.success`)
	base.Success(ctx, data)
}

func (c *Controller) UpdateOrderController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.update.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	var req UpdateOrderControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	addressID, err := uuid.Parse(req.AddressID)
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

	if err := c.svc.UpdateOrderService(ctx.Request.Context(), orderID, &UpdateOrderServiceRequest{
		PaymentID:          paymentID,
		AddressID:          addressID,
		Status:             req.Status,
		ShippingTrackingNo: req.ShippingTrackingNo,
		TotalAmount:        req.TotalAmount,
		DiscountAmount:     req.DiscountAmount,
		NetAmount:          req.NetAmount,
	}, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteOrderController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`orders.ctl.delete.start`)

	orderID, ok := c.parseOrderID(ctx)
	if !ok {
		return
	}

	requesterID, hasRequester := auth.GetMemberID(ctx)
	isAdmin := auth.GetIsAdmin(ctx)
	if !isAdmin && !hasRequester {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	if err := c.svc.DeleteOrderService(ctx.Request.Context(), orderID, requesterID, isAdmin); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`orders.ctl.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseOrderID(ctx *gin.Context) (uuid.UUID, bool) {
	var uri OrderURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, false
	}

	orderID, err := uuid.Parse(uri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, false
	}

	return orderID, true
}
