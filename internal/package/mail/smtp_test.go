package mail

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

func TestVerificationEmailTemplateRenders(t *testing.T) {
	tmpl, err := template.ParseFS(templateFS, "templates/email_verification.html")
	if err != nil {
		t.Fatalf("parse verification email template: %v", err)
	}

	data := VerificationEmailData{
		VerifyURL: "https://example.com/api/v1/auth/verify-email?token=test-token",
		Year:      2026,
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		t.Fatalf("render verification email template: %v", err)
	}

	for _, want := range []string{
		`href="` + data.VerifyURL + `"`,
		data.VerifyURL,
		"© 2026 PasarDesa",
	} {
		if !strings.Contains(body.String(), want) {
			t.Errorf("rendered template does not contain %q", want)
		}
	}
}
