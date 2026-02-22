package contact

import (
	"strings"

	"phakram/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateReplyPublicControllerRequest struct {
	Message string `json:"message"`
}

func (c *Controller) ListRepliesPublicController(ctx *gin.Context) {
	contactMessageID := strings.TrimSpace(ctx.Param("id"))
	chatToken := strings.TrimSpace(ctx.Query("token"))

	if contactMessageID == "" || chatToken == "" {
		base.BadRequest(ctx, "ข้อมูลไม่ครบถ้วน", nil)
		return
	}

	data, err := c.svc.ListRepliesPublic(ctx.Request.Context(), contactMessageID, chatToken)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, data)
}

func (c *Controller) CreateReplyPublicController(ctx *gin.Context) {
	contactMessageID := strings.TrimSpace(ctx.Param("id"))
	chatToken := strings.TrimSpace(ctx.Query("token"))

	if contactMessageID == "" || chatToken == "" {
		base.BadRequest(ctx, "ข้อมูลไม่ครบถ้วน", nil)
		return
	}

	var req CreateReplyPublicControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "ข้อมูลไม่ถูกต้อง", nil)
		return
	}

	replyMessage := strings.TrimSpace(req.Message)
	if replyMessage == "" {
		base.BadRequest(ctx, "กรุณาระบุข้อความ", nil)
		return
	}
	if len(replyMessage) > 3000 {
		base.BadRequest(ctx, "ข้อความยาวเกินกำหนด", nil)
		return
	}

	result, err := c.svc.CreateReplyPublic(ctx.Request.Context(), contactMessageID, chatToken, replyMessage)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, result)
}
