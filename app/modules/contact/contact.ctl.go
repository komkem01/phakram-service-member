package contact

import (
	"net/mail"
	"strings"

	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
)

type SubmitContactController struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (c *Controller) SubmitController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent("contact.ctl.submit.start")

	var req SubmitContactController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Subject = strings.TrimSpace(req.Subject)
	req.Message = strings.TrimSpace(req.Message)

	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
		base.BadRequest(ctx, "กรุณากรอกข้อมูลให้ครบถ้วน", nil)
		return
	}

	if len(req.Name) > 120 || len(req.Email) > 200 || len(req.Subject) > 200 || len(req.Message) > 5000 {
		base.BadRequest(ctx, "ข้อมูลยาวเกินกำหนด", nil)
		return
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		base.BadRequest(ctx, "รูปแบบอีเมลไม่ถูกต้อง", nil)
		return
	}

	result, err := c.svc.Submit(ctx.Request.Context(), &SubmitContactService{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent("contact.ctl.submit.success")
	base.Success(ctx, result)
}
