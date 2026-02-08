package product_files

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateProductFileControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateProductFileController struct {
	ProductID string `json:"product_id"`
	FileID    string `json:"file_id"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateProductFileControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_files.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateProductFileController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_files.ctl.update.request_body`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	var productID uuid.UUID
	if req.ProductID != "" {
		productID, err = uuid.Parse(req.ProductID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var fileID uuid.UUID
	if req.FileID != "" {
		fileID, err = uuid.Parse(req.FileID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateProductFileService{
		ProductID: productID,
		FileID:    fileID,
		MemberID:  memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`product_files.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) ProductFilesUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
