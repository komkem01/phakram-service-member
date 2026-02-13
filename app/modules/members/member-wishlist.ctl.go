package members

import (
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberWishlistURIRequest struct {
	MemberID   string `uri:"id"`
	WishlistID string `uri:"wishlist_id"`
}

type ListMemberWishlistControllerRequest struct {
	base.RequestPaginate
}

type CreateMemberWishlistControllerRequest struct {
	ProductID       string `json:"product_id"`
	Quantity        int    `json:"quantity"`
	PricePerUnit    string `json:"price_per_unit"`
	TotalItemAmount string `json:"total_item_amount"`
}

type UpdateMemberWishlistControllerRequest = CreateMemberWishlistControllerRequest

func (c *Controller) ListMemberWishlistController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.wishlist.list.start`)

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req ListMemberWishlistControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListMemberWishlistService(ctx.Request.Context(), &entitiesdto.ListMemberWishlistRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.wishlist.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) CreateMemberWishlistController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.wishlist.create.start`)

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req CreateMemberWishlistControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.CreateMemberWishlistService(ctx.Request.Context(), memberID, &CreateMemberWishlistServiceRequest{
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
		ActionBy:        actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.wishlist.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) InfoMemberWishlistController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.wishlist.info.start`)

	memberID, wishlistID, ok := c.parseMemberWishlistURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	data, err := c.svc.InfoMemberWishlistService(ctx.Request.Context(), memberID, wishlistID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.wishlist.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) UpdateMemberWishlistController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.wishlist.update.start`)

	memberID, wishlistID, ok := c.parseMemberWishlistURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	var req UpdateMemberWishlistControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.UpdateMemberWishlistService(ctx.Request.Context(), memberID, wishlistID, &UpdateMemberWishlistServiceRequest{
		ProductID:       productID,
		Quantity:        req.Quantity,
		PricePerUnit:    req.PricePerUnit,
		TotalItemAmount: req.TotalItemAmount,
		ActionBy:        actionBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.wishlist.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteMemberWishlistController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.wishlist.delete.start`)

	memberID, wishlistID, ok := c.parseMemberWishlistURI(ctx)
	if !ok {
		return
	}

	if !c.ensureAdminOrSelf(ctx, memberID) {
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.DeleteMemberWishlistService(ctx.Request.Context(), memberID, wishlistID, actionBy); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.wishlist.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseMemberWishlistURI(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	var uri MemberWishlistURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	memberID, err := uuid.Parse(uri.MemberID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	wishlistID, err := uuid.Parse(uri.WishlistID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return memberID, wishlistID, true
}
