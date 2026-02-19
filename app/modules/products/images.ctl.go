package products

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type productImageControllerURI struct {
	ID string `uri:"id"`
}

type deleteProductImageControllerURI struct {
	ID      string `uri:"id"`
	ImageID string `uri:"image_id"`
}

type uploadProductImageControllerRequest struct {
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	FileSize   int64  `json:"file_size"`
	FileBase64 string `json:"file_base64"`
}

func (c *Controller) ListProductImagesController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`products.ctl.images.list.start`)

	var uri productImageControllerURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(uri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, err := c.svc.ListProductImagesService(ctx.Request.Context(), productID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`products.ctl.images.list.success`)
	base.Success(ctx, data)
}

func (c *Controller) UploadProductImageController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`products.ctl.images.upload.start`)

	var uri productImageControllerURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(uri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req uploadProductImageControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	serviceReq := &UploadProductImageServiceRequest{
		FileName:   req.FileName,
		FileType:   req.FileType,
		FileSize:   req.FileSize,
		FileBase64: req.FileBase64,
	}

	if err := c.svc.normalizeProductImageInput(serviceReq); err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	item, err := c.svc.UploadProductImageService(ctx.Request.Context(), productID, serviceReq)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`products.ctl.images.upload.success`)
	base.Success(ctx, item)
}

func (c *Controller) DeleteProductImageController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`products.ctl.images.delete.start`)

	var uri deleteProductImageControllerURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	productID, err := uuid.Parse(uri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	imageID, err := uuid.Parse(uri.ImageID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.DeleteProductImageService(ctx.Request.Context(), productID, imageID); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`products.ctl.images.delete.success`)
	base.Success(ctx, nil)
}
