package router

import (
	"go-starterkit/internal/config"
	"go-starterkit/internal/middleware"
	"go-starterkit/internal/modules/auth"
	user "go-starterkit/internal/modules/users"
	"go-starterkit/internal/package/jwt"
	"go-starterkit/internal/package/logger"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// wire set
var RouterSet = wire.NewSet(
	NewAppHandler,
	middleware.RateLimitSet,
)

type AppHandler struct {
	config *config.Config

	// Add Handler In Here
	log         logger.Logger
	jwtSvc      jwt.JWTService
	rateLimiter *middleware.RateLimitMiddleware
	UserHandler *user.HandlerUser
	RouteAuth   *auth.RouterAuth
	Enforcer    *casbin.Enforcer
}

func NewAppHandler(
	cfg *config.Config,
	log logger.Logger,
	jwtSvc jwt.JWTService,
	rateLimiter *middleware.RateLimitMiddleware,
	userH *user.HandlerUser,
	RouteAuth *auth.RouterAuth,
	Enforcer *casbin.Enforcer,
) *AppHandler {
	return &AppHandler{
		config:      cfg,
		jwtSvc:      jwtSvc,
		rateLimiter: rateLimiter,
		UserHandler: userH,
		RouteAuth:   RouteAuth,
		Enforcer:    Enforcer,
	}
}

func (app *AppHandler) RegisterRoutes(r *gin.Engine) {

	r.Use(middleware.CORS())
	r.Use(middleware.LoggerMiddleware(app.log))

	r.Use(gin.Recovery())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Base API V1 Route Group
	apiV1 := r.Group("/api/v1")
	apiV1.Use(app.rateLimiter.LimitByIP("global"))
	{
		app.registerPublicRoutes(apiV1)
		app.registerProtectedRoutes(apiV1)
	}

	// Custom 404 Handler (Biar response JSON-nya seragam)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Route not found",
		})
	})
}
