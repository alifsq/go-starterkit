package bootstrap

import (
	"go-starterkit/internal/config"
	"go-starterkit/internal/package/logger"
	"go-starterkit/internal/package/mail"
)

func NewZapLogger(cfg *config.Config) (logger.Logger, error) {
	return logger.NewLogger(cfg.AppEnv)
}

func NewMailerConfig(cfg *config.Config) mail.MailerConfig {
	return mail.MailerConfig{
		AppBaseURL:   cfg.AppBaseURL,
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		SMTPUsername: cfg.SMTPUsername,
		SMTPPassword: cfg.SMTPPassword,
		SMTPFrom:     cfg.SMTPFrom,
	}
}
