package member_files

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadMemberFileRequestUri struct {
	MemberID string `uri:"member_id" binding:"required"`
}

type UploadMemberFileResponse struct {
	FileID       string `json:"file_id"`
	MemberFileID string `json:"member_file_id"`
	FilePath     string `json:"file_path"`
}

func (c *Controller) UploadMemberFileController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`member_files.ctl.upload.start`)

	var reqUri UploadMemberFileRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`member_files.ctl.upload.request_uri`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	refID, err := uuid.Parse(reqUri.MemberID)
	if err != nil {
		log.With(slog.Any(`uri`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	originalName := fileHeader.Filename
	ext := filepath.Ext(originalName)
	fileName := uuid.New().String() + ext
	fileDir := filepath.Join("uploads", "members", refID.String())
	if err := os.MkdirAll(fileDir, 0o755); err != nil {
		log.Errf(`internal: %s`, err)
		base.InternalServerError(ctx, i18n.InternalServerError, nil)
		return
	}

	filePath := filepath.Join(fileDir, fileName)
	if err := ctx.SaveUploadedFile(fileHeader, filePath); err != nil {
		log.Errf(`internal: %s`, err)
		base.InternalServerError(ctx, i18n.InternalServerError, nil)
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}
	fileSize := strconv.FormatInt(fileHeader.Size, 10)

	resp, err := c.svc.UploadMemberFileService(ctx.Request.Context(), &UploadMemberFileService{
		MemberID: refID,
		FileName: originalName,
		FilePath: filepath.ToSlash(filePath),
		FileType: fileType,
		FileSize: fileSize,
		MemberBy: memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`member_files.ctl.upload.success`)
	base.Success(ctx, resp)
}
