package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReadRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		cookie    string
		body      string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "reads cookie",
			cookie:    "cookie-token",
			body:      `{"refresh_token":"body-token"}`,
			wantToken: "cookie-token",
		},
		{
			name:      "reads JSON body when cookie is absent",
			body:      `{"refresh_token":"body-token"}`,
			wantToken: "body-token",
		},
		{
			name: "returns empty token when cookie and body are absent",
		},
		{
			name:    "rejects malformed JSON",
			body:    `{"refresh_token":`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.cookie})
			}

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = req

			got, err := readRefreshToken(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("readRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.wantToken {
				t.Fatalf("readRefreshToken() = %q, want %q", got, tt.wantToken)
			}
		})
	}
}
