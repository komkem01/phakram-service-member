package orders

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListOrderControllerRequest struct {
	base.RequestPaginate
	Status string `form:"status"`
}

type ListOrderByMemberRequestUri struct {
	ID string `uri:"id"`
}

type ListOrderControllerResponses struct {
	ID             uuid.UUID       `json:"id"`
	OrderNo        string          `json:"order_no"`
	MemberID       uuid.UUID       `json:"member_id"`
	PaymentID      uuid.UUID       `json:"payment_id"`
	AddressID      uuid.UUID       `json:"address_id"`
	Status         string          `json:"status"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	NetAmount      decimal.Decimal `json:"net_amount"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

func (c *Controller) OrdersList(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var req ListOrderControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`orders.ctl.list.request`)

	memberID := uuid.Nil
	isAdmin := authmod.GetIsAdmin(ctx)
	if !isAdmin {
		var ok bool
		memberID, ok = authmod.GetMemberID(ctx)
		if !ok {
			base.Unauthorized(ctx, i18n.Unauthorized, nil)
			return
		}
	}

	data, page, err := c.svc.ListService(ctx, &ListOrderServiceRequest{
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
	span.AddEvent(`orders.ctl.list.callsvc`)

	var resp []*ListOrderControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}

func (c *Controller) OrdersListByMember(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	if !authmod.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var reqUri ListOrderByMemberRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	memberID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req ListOrderControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListService(ctx, &ListOrderServiceRequest{
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
	span.AddEvent(`orders.ctl.list_by_member.callsvc`)

	var resp []*ListOrderControllerResponses
	if err := utils.CopyNTimeToUnix(&resp, data); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.InternalServerError(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, resp, page)
}
