package mail

import (
	"context"

	"go-starterkit/internal/package/logger"
)

type Mailer interface {
	SendVerificationEmail(ctx context.Context, toEmail, token string) error
	SendResetPasswordEmail(ctx context.Context, toEmail, token string) error
}

type MailerConfig struct {
	AppBaseURL   string
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

type smtpMailer struct {
	cfg MailerConfig
	log logger.Logger
}

func NewMailer(cfg MailerConfig, log logger.Logger) Mailer {
	return &smtpMailer{
		cfg: cfg,
		log: log,
	}
}
