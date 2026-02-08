package storages

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateStorageControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateStorageController struct {
	RefID         string `json:"ref_id"`
	FileName      string `json:"file_name"`
	FilePath      string `json:"file_path"`
	FileType      string `json:"file_type"`
	FileSize      string `json:"file_size"`
	RelatedEntity string `json:"related_entity"`
	UploadedBy    string `json:"uploaded_by"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateStorageControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`storages.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateStorageController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`storages.ctl.update.request_body`)

	var refID uuid.UUID
	if req.RefID != "" {
		refID, err = uuid.Parse(req.RefID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}
	var uploadedBy uuid.UUID
	if req.UploadedBy != "" {
		uploadedBy, err = uuid.Parse(req.UploadedBy)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
	}

	if err := c.svc.UpdateService(ctx, id, &UpdateStorageService{
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
	span.AddEvent(`storages.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) StoragesUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
