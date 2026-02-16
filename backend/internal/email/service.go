package email

import (
	"fmt"
	"net/smtp"

	"github.com/monachy/projek/internal/common"
)

type Service interface {
	Send(to, subject, body string) error
	SendTemplate(to, subject string, data interface{}) error
}

type SMTPService struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

func NewSMTPService(cfg *common.Config) *SMTPService {
	return &SMTPService{
		host:     cfg.Email.SMTPHost,
		port:     cfg.Email.SMTPPort,
		username: cfg.Email.SMTPUser,
		password: cfg.Email.SMTPPassword,
		from:     cfg.Email.FromEmail,
		fromName: cfg.Email.FromName,
	}
}

func (s *SMTPService) Send(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.fromName, s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

func (s *SMTPService) SendTemplate(to, subject string, data interface{}) error {
	return s.Send(to, subject, fmt.Sprintf("%v", data))
}

func SendInvitationEmail(svc Service, to, inviterName, communityName, token string) error {
	subject := "You're invited to join " + communityName
	body := fmt.Sprintf(`
Hello,

%s has invited you to join %s on Projek.

Click the link below to accept the invitation:
https://app.projek.local/invite/%s

This invitation will expire in 7 days.
	`, inviterName, communityName, token)

	return svc.Send(to, subject, body)
}

func SendWelcomeEmail(svc Service, to, firstName, communityName string) error {
	subject := "Welcome to " + communityName
	body := fmt.Sprintf(`
Hello %s,

Welcome to %s! We're excited to have you on board.

Get started by creating your first project and inviting your team members.

Best regards,
The Projek Team
	`, firstName, communityName)

	return svc.Send(to, subject, body)
}
