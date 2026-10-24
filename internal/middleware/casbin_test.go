package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-starterkit/internal/middleware"
	jwtservice "go-starterkit/internal/package/jwt"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/gin-gonic/gin"
)

func TestAuthAndCasbinAuthorizeUserMeByRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authModel, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act
`)
	if err != nil {
		t.Fatal(err)
	}

	enforcer, err := casbin.NewEnforcer(authModel)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := enforcer.AddPolicy("user", "/api/v1/user/me", "GET"); err != nil {
		t.Fatal(err)
	}

	jwtSvc := jwtservice.NewJWTService("test-secret")
	token, err := jwtSvc.GenerateAccessToken(42, "user", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.Use(middleware.Auth(jwtSvc), middleware.CasbinMiddleware(enforcer))
	router.GET("/api/v1/user/me", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	var body struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(strings.NewReader(recorder.Body.String())).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.UserID != 42 {
		t.Fatalf("expected user ID 42, got %d", body.UserID)
	}
}
