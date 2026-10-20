package auth

import "github.com/gin-gonic/gin"

func (h *HandlerAuth) RegisterPublicRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.Refresh)
		authGroup.POST("/logout", h.Logout)

		// authGroup.POST("/resend-verification", h.ResendVerification)

		authGroup.POST("/verify-email", h.VerifyEmail)
		authGroup.POST("/forgot-password", h.ForgotPassword)
		authGroup.POST("/reset-password", h.ResetPassword)
	}
}
func (h *HandlerAuth) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		// authGroup.GET("/me", h.Me)
		authGroup.POST("/revoke-all", h.RevokeAll)
	}
}
