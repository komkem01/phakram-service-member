package contact

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type contactMessageRecord struct {
	bun.BaseModel `bun:"table:contact_messages"`

	ID         uuid.UUID  `bun:"id,pk,type:uuid"`
	Name       string     `bun:"name,notnull"`
	Email      string     `bun:"email,notnull"`
	Subject    string     `bun:"subject,notnull"`
	Message    string     `bun:"message,notnull"`
	SendStatus string     `bun:"send_status,notnull"`
	IsRead     bool       `bun:"is_read,notnull"`
	SendError  string     `bun:"send_error"`
	SentAt     *time.Time `bun:"sent_at"`
	ReadAt     *time.Time `bun:"read_at"`
	CreatedAt  time.Time  `bun:"created_at,notnull"`
	UpdatedAt  time.Time  `bun:"updated_at,notnull"`
}

type SubmitContactService struct {
	Name    string
	Email   string
	Subject string
	Message string
}

type SubmitContactResult struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (s *Service) Submit(ctx context.Context, req *SubmitContactService) (*SubmitContactResult, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent("contact.svc.submit.start")

	now := time.Now()
	id := uuid.New()

	record := &contactMessageRecord{
		ID:         id,
		Name:       req.Name,
		Email:      req.Email,
		Subject:    req.Subject,
		Message:    req.Message,
		SendStatus: "pending",
		IsRead:     false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if _, err := s.bunDB.DB().NewInsert().Model(record).Exec(ctx); err != nil {
		return nil, err
	}

	if err := s.sendMail(req); err != nil {
		_, _ = s.bunDB.DB().NewUpdate().
			Model((*contactMessageRecord)(nil)).
			Set("send_status = ?", "failed").
			Set("send_error = ?", err.Error()).
			Set("updated_at = ?", time.Now()).
			Where("id = ?", id).
			Exec(ctx)
		return nil, err
	}

	sentAt := time.Now()
	_, _ = s.bunDB.DB().NewUpdate().
		Model((*contactMessageRecord)(nil)).
		Set("send_status = ?", "sent").
		Set("send_error = NULL").
		Set("sent_at = ?", sentAt).
		Set("updated_at = ?", sentAt).
		Where("id = ?", id).
		Exec(ctx)

	span.AddEvent("contact.svc.submit.success")
	return &SubmitContactResult{
		ID:      id.String(),
		Message: "ส่งข้อความเรียบร้อยแล้ว",
	}, nil
}

func (s *Service) sendMail(req *SubmitContactService) error {
	if s.conf == nil {
		return fmt.Errorf("contact config is not set")
	}

	host := strings.TrimSpace(s.conf.Mail.Host)
	port := s.conf.Mail.Port
	username := strings.TrimSpace(s.conf.Mail.Username)
	password := strings.TrimSpace(s.conf.Mail.Password)
	from := strings.TrimSpace(s.conf.Mail.From)
	to := strings.TrimSpace(s.conf.RecipientEmail)

	if to == "" {
		to = "komkem.contact@gmail.com"
	}
	if host == "" || port <= 0 || from == "" {
		return fmt.Errorf("smtp is not configured")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.Auth(nil)
	if username != "" && password != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	safeSubject := strings.ReplaceAll(req.Subject, "\n", " ")
	safeSubject = strings.ReplaceAll(safeSubject, "\r", " ")

	body := fmt.Sprintf(
		"ได้รับข้อความจากแบบฟอร์มติดต่อ\n\nชื่อ: %s\nอีเมล: %s\nหัวข้อ: %s\n\nข้อความ:\n%s\n",
		req.Name,
		req.Email,
		req.Subject,
		req.Message,
	)

	message := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Reply-To: %s", req.Email),
		fmt.Sprintf("Subject: [Phakram Contact] %s", safeSubject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
}
