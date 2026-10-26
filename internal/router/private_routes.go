package router

import (
	"go-starterkit/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (app *AppHandler) registerProtectedRoutes(rg *gin.RouterGroup) {
	// Buat group yang dilindungi Middleware JWT Auth
	protected := rg.Group("")
	protected.Use(
		middleware.Auth(app.jwtSvc),
		middleware.CasbinMiddleware(app.Enforcer),
	)
	{
		// Modul Auth Protected Routes (/revoke-all)
		app.RouteAuth.RegisterProtectedRoutes(protected)

		// Modul Users Routes
		app.UserHandler.RoutesUser(protected)

		// Tambahkan route protected modul lain di sini
	}
}
