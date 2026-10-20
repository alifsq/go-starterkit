package mail

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

//go:embed templates/*.html
var templateFS embed.FS

type VerificationEmailData struct {
	VerifyURL string
	Year      int
}

type PasswordResetData struct {
	ResetURL string
}

func (m *smtpMailer) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	m.log.Info("Mengirim email verifikasi", zap.String("email", toEmail))

	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", m.cfg.AppBaseURL, token)

	tmpl, err := template.ParseFS(templateFS, "templates/email_verification.html")
	if err != nil {
		m.log.Error("Gagal parse template email verifikasi", zap.Error(err))
		return err
	}

	var bodyBuffer bytes.Buffer
	data := VerificationEmailData{
		VerifyURL: verifyURL,
		Year:      time.Now().Year(),
	}
	if err := tmpl.Execute(&bodyBuffer, data); err != nil {
		m.log.Error("Gagal render template email verifikasi", zap.Error(err))
		return err
	}

	return m.sendMail(toEmail, "Verifikasi Email Akun Anda", bodyBuffer.String())
}

func (m *smtpMailer) SendResetPasswordEmail(ctx context.Context, toEmail, token string) error {
	m.log.Info("Mengirim email reset password", zap.String("email", toEmail))

	resetURL := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", m.cfg.AppBaseURL, token)

	tmpl, err := template.ParseFS(templateFS, "templates/password_reset.html")
	if err != nil {
		m.log.Error("Gagal parse template email reset password", zap.Error(err))
		return err
	}

	var bodyBuffer bytes.Buffer
	data := PasswordResetData{ResetURL: resetURL}
	if err := tmpl.Execute(&bodyBuffer, data); err != nil {
		m.log.Error("Gagal render template email reset password", zap.Error(err))
		return err
	}

	return m.sendMail(toEmail, "Permintaan Reset Password", bodyBuffer.String())
}

func (m *smtpMailer) sendMail(toEmail, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.cfg.SMTPFrom)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	port, err := strconv.Atoi(m.cfg.SMTPPort)
	if err != nil {
		port = 587 // Fallback default port
	}

	dialer := gomail.NewDialer(m.cfg.SMTPHost, port, m.cfg.SMTPUsername, m.cfg.SMTPPassword)

	if err := dialer.DialAndSend(msg); err != nil {
		m.log.Error("Gagal mengirim email via SMTP dialer", zap.Error(err), zap.String("to", toEmail))
		return err
	}

	m.log.Info("Email berhasil dikirim", zap.String("to", toEmail))
	return nil
}
