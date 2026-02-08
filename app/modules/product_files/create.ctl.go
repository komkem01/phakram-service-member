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

type CreateProductFileController struct {
	ProductID string `json:"product_id"`
	FileID    string `json:"file_id"`
}

func (c *Controller) CreateProductFileController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`product_files.ctl.create.start`)

	var req CreateProductFileController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`product_files.ctl.create.request`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	fileID, err := uuid.Parse(req.FileID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateProductFileService(ctx.Request.Context(), &CreateProductFileService{
		ProductID: productID,
		FileID:    fileID,
		MemberID:  memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`product_files.ctl.create.success`)
	base.Success(ctx, nil)
}
