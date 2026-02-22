package contact

import (
	"strings"

	"phakram/app/modules/auth"
	"phakram/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateReplyControllerRequest struct {
	Message string `json:"message"`
}

func (c *Controller) ListRepliesController(ctx *gin.Context) {
	contactMessageID := strings.TrimSpace(ctx.Param("id"))
	if contactMessageID == "" {
		base.BadRequest(ctx, "ไม่พบรหัสข้อความติดต่อ", nil)
		return
	}

	data, err := c.svc.ListReplies(ctx.Request.Context(), contactMessageID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, data)
}

func (c *Controller) CreateReplyController(ctx *gin.Context) {
	if !auth.GetActorIsAdmin(ctx) {
		base.Forbidden(ctx, "ไม่มีสิทธิ์เข้าถึง", nil)
		return
	}

	contactMessageID := strings.TrimSpace(ctx.Param("id"))
	if contactMessageID == "" {
		base.BadRequest(ctx, "ไม่พบรหัสข้อความติดต่อ", nil)
		return
	}

	var req CreateReplyControllerRequest
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

	result, err := c.svc.CreateReply(ctx.Request.Context(), &CreateReplyServiceRequest{
		ContactMessageID: contactMessageID,
		SenderRole:       "admin",
		SenderName:       "ทีมงาน",
		Message:          replyMessage,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	base.Success(ctx, result)
}
