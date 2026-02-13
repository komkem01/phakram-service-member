package storages

import (
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateStorageControllerRequest struct {
	RefID         string `json:"ref_id"`
	FileName      string `json:"file_name"`
	FilePath      string `json:"file_path"`
	FileSize      int64  `json:"file_size"`
	FileType      string `json:"file_type"`
	RelatedEntity string `json:"related_entity"`
	UploadedBy    string `json:"uploaded_by"`
	IsActive      *bool  `json:"is_active"`
}

func (c *Controller) CreateStorageController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`storages.ctl.create.start`)

	var req CreateStorageControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	refID, err := uuid.Parse(req.RefID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	uploadedBy, err := uuid.Parse(req.UploadedBy)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	if err := c.svc.CreateService(ctx.Request.Context(), &CreateStorageServiceRequest{
		RefID:         refID,
		FileName:      req.FileName,
		FilePath:      req.FilePath,
		FileSize:      req.FileSize,
		FileType:      req.FileType,
		RelatedEntity: req.RelatedEntity,
		UploadedBy:    uploadedBy,
		IsActive:      req.IsActive,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`storages.ctl.create.success`)
	base.Success(ctx, nil)
}
