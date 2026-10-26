package auth

import (
	"fmt"
	"go-starterkit/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RouterAuth struct {
	authHandler *HandlerAuth
	rateLimiter *middleware.RateLimitMiddleware // <-- Tambahkan ini
}

// Constructor yang akan dirakit otomatis oleh Google Wire
func NewRouterAuth(authHandler *HandlerAuth, rateLimiter *middleware.RateLimitMiddleware) *RouterAuth {
	return &RouterAuth{
		authHandler: authHandler,
		rateLimiter: rateLimiter,
	}
}

func (h *RouterAuth) RegisterPublicRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.rateLimiter.LimitByIP("register_ip"), h.authHandler.Register)
		authGroup.POST("/login", h.rateLimiter.LimitByIP("login_ip"), h.authHandler.Login)
		authGroup.POST("/refresh", h.authHandler.Refresh)
		authGroup.POST("/logout", h.authHandler.Logout)

		authGroup.POST("/verify-email", h.rateLimiter.LimitByIP("verify_email_ip"), h.authHandler.VerifyEmail)
		authGroup.POST("/resend-verification", h.rateLimiter.LimitByUser("resend_verification", func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("%v", userID)
			}
			return ""
		}), h.authHandler.ResendVerificationEmail)
		authGroup.POST("/forgot-password", h.rateLimiter.LimitByIP("email_send_ip"), h.authHandler.ForgotPassword)
		authGroup.POST("/reset-password", h.rateLimiter.LimitByIP("password_reset_ip"), h.authHandler.ResetPassword)
	}
}
func (h *RouterAuth) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		// authGroup.GET("/me", h.Me)
		authGroup.POST("/revoke-all", h.rateLimiter.LimitByUser("user_action_std", func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("%v", userID)
			}
			return ""
		}), h.authHandler.RevokeAll)
	}
}
