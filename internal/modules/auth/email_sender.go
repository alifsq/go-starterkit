package auth

import (
	"context"
	"fmt"
	"net/smtp"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, toEmail, token string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, token string) error
}

type noopEmailSender struct{}

func (n *noopEmailSender) SendVerificationEmail(_ context.Context, _, _ string) error  { return nil }
func (n *noopEmailSender) SendPasswordResetEmail(_ context.Context, _, _ string) error { return nil }

type smtpEmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	baseURL  string
}

func NewEmailSender() EmailSender {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	portStr := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if host == "" || portStr == "" {
		return &noopEmailSender{}
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return &noopEmailSender{}
	}
	return &smtpEmailSender{
		host:     host,
		port:     port,
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
		baseURL:  strings.TrimRight(strings.TrimSpace(os.Getenv("APP_FRONTEND_URL")), "/"),
	}
}

func (s *smtpEmailSender) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	if s.from == "" || toEmail == "" {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	msg := buildEmailMessage("Verify your email", toEmail, fmt.Sprintf("Welcome! Please verify your email using this link: %s", s.emailLink("verify-email", token)))
	return s.sendMail(toEmail, msg)
}

func (s *smtpEmailSender) SendPasswordResetEmail(ctx context.Context, toEmail, token string) error {
	if s.from == "" || toEmail == "" {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	msg := buildEmailMessage("Reset your password", toEmail, fmt.Sprintf("Use this link to reset your password: %s", s.emailLink("reset-password", token)))
	return s.sendMail(toEmail, msg)
}

func (s *smtpEmailSender) emailLink(path, token string) string {
	if s.baseURL == "" {
		return path + "?token=" + url.QueryEscape(token)
	}
	return s.baseURL + "/" + strings.TrimLeft(path, "/") + "?token=" + url.QueryEscape(token)
}

func buildEmailMessage(subject, toEmail, bodyHTML string) []byte {
	return []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		bodyHTML + "\r\n")
}
func (s *smtpEmailSender) sendMail(toEmail string, content []byte) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(s.host+":"+strconv.Itoa(s.port), auth, s.from, []string{toEmail}, content)
}
