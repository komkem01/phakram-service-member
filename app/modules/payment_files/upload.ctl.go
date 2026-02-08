package payment_files

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

type UploadPaymentFileRequestUri struct {
	PaymentID string `uri:"payment_id" binding:"required"`
}

type UploadPaymentFileResponse struct {
	FileID        string `json:"file_id"`
	PaymentFileID string `json:"payment_file_id"`
	FilePath      string `json:"file_path"`
}

func (c *Controller) UploadPaymentFileController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`payment_files.ctl.upload.start`)

	var reqUri UploadPaymentFileRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`payment_files.ctl.upload.request_uri`)

	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	paymentID, err := uuid.Parse(reqUri.PaymentID)
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
	fileDir := filepath.Join("uploads", "payments", paymentID.String())
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

	resp, err := c.svc.UploadPaymentFileService(ctx.Request.Context(), &UploadPaymentFileService{
		PaymentID: paymentID,
		FileName:  originalName,
		FilePath:  filepath.ToSlash(filePath),
		FileType:  fileType,
		FileSize:  fileSize,
		MemberID:  memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`payment_files.ctl.upload.success`)
	base.Success(ctx, resp)
}
