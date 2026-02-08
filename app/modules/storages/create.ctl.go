package storages

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateStorageController struct {
	RefID         string `json:"ref_id"`
	FileName      string `json:"file_name"`
	FilePath      string `json:"file_path"`
	FileType      string `json:"file_type"`
	FileSize      string `json:"file_size"`
	RelatedEntity string `json:"related_entity"`
	UploadedBy    string `json:"uploaded_by"`
}

func (c *Controller) CreateStorageController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`storages.ctl.create.start`)

	var req CreateStorageController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`storages.ctl.create.request`)

	refID, err := uuid.Parse(req.RefID)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	uploadedBy, err := uuid.Parse(req.UploadedBy)
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateStorageService(ctx.Request.Context(), &CreateStorageService{
		RefID:         refID,
		FileName:      req.FileName,
		FilePath:      req.FilePath,
		FileType:      req.FileType,
		FileSize:      req.FileSize,
		RelatedEntity: req.RelatedEntity,
		UploadedBy:    uploadedBy,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`storages.ctl.create.success`)
	base.Success(ctx, nil)
}
