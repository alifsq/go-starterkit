package auth

import "github.com/gin-gonic/gin"

func (h *HandlerAuth) RegisterPublicRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.Refresh) // Membaca HttpOnly Cookie
		authGroup.POST("/logout", h.Logout)
	}
}
func (h *HandlerAuth) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/revoke-all", h.RevokeAll)
	}
}
