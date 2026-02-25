package reviews

import (
	"strings"

	"phakram/app/modules/auth"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateReviewControllerRequest struct {
	OrderItemID string   `json:"order_item_id"`
	Rating      int      `json:"rating"`
	Comment     string   `json:"comment"`
	ImageURLs   []string `json:"image_urls"`
}

type UpdateReviewControllerRequest struct {
	Rating    int      `json:"rating"`
	Comment   string   `json:"comment"`
	ImageURLs []string `json:"image_urls"`
}

type UpdateReviewVisibilityControllerRequest struct {
	IsVisible bool `json:"is_visible"`
}

type ProductReviewURIRequest struct {
	ProductID string `uri:"id"`
}

type ReviewURIRequest struct {
	ReviewID string `uri:"id"`
}

type EligibleReviewQueryRequest struct {
	ProductID string `form:"product_id"`
}

type ListAdminReviewControllerRequest struct {
	base.RequestPaginate
	ProductID string `form:"product_id"`
	IsVisible *bool  `form:"is_visible"`
	HasImages *bool  `form:"has_images"`
	Rating    *int   `form:"rating"`
}

type ListProductReviewsControllerRequest struct {
	base.RequestPaginate
	HasImages *bool `form:"has_images"`
	Rating    *int  `form:"rating"`
}

func (c *Controller) ListProductPublicController(ctx *gin.Context) {
	var uriReq ProductReviewURIRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(strings.TrimSpace(uriReq.ProductID))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req ListProductReviewsControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListPublicByProduct(ctx.Request.Context(), productID, &ListProductReviewsServiceRequest{
		RequestPaginate: req.RequestPaginate,
		HasImages:       req.HasImages,
		Rating:          req.Rating,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Paginate(ctx, data, page)
}

func (c *Controller) ListEligibleController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var req EligibleReviewQueryRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(strings.TrimSpace(req.ProductID))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.ListEligibleByProduct(ctx.Request.Context(), memberID, productID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, data)
}

func (c *Controller) CreateController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var req CreateReviewControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	orderItemID, err := uuid.Parse(strings.TrimSpace(req.OrderItemID))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.Create(ctx.Request.Context(), &CreateReviewServiceRequest{
		MemberID:    memberID,
		OrderItemID: orderItemID,
		Rating:      req.Rating,
		Comment:     req.Comment,
		ImageURLs:   req.ImageURLs,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var uriReq ReviewURIRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	reviewID, err := uuid.Parse(strings.TrimSpace(uriReq.ReviewID))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateReviewControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.Update(ctx.Request.Context(), &UpdateReviewServiceRequest{
		ReviewID:  reviewID,
		MemberID:  memberID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		ImageURLs: req.ImageURLs,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	memberID, ok := auth.GetMemberID(ctx)
	if !ok {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var uriReq ReviewURIRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	reviewID, err := uuid.Parse(strings.TrimSpace(uriReq.ReviewID))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.Delete(ctx.Request.Context(), reviewID, memberID); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}

func (c *Controller) ListAdminController(ctx *gin.Context) {
	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var req ListAdminReviewControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var productID *uuid.UUID
	if strings.TrimSpace(req.ProductID) != "" {
		parsedProductID, err := uuid.Parse(strings.TrimSpace(req.ProductID))
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		productID = &parsedProductID
	}

	data, page, err := c.svc.ListAdmin(ctx.Request.Context(), &ListAdminReviewsServiceRequest{
		RequestPaginate: req.RequestPaginate,
		ProductID:       productID,
		IsVisible:       req.IsVisible,
		HasImages:       req.HasImages,
		Rating:          req.Rating,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Paginate(ctx, data, page)
}

func (c *Controller) UpdateVisibilityController(ctx *gin.Context) {
	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	var uriReq ReviewURIRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	reviewID, err := uuid.Parse(strings.TrimSpace(uriReq.ReviewID))
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateReviewVisibilityControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.UpdateVisibility(ctx.Request.Context(), reviewID, req.IsVisible); err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, nil)
}
