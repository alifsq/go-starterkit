package router

import "github.com/gin-gonic/gin"

func (app *AppHandler) registerPublicRoutes(rg *gin.RouterGroup) {
	// Modul Auth Public Routes (/register, /login, /refresh, /logout)
	app.RouteAuth.RegisterPublicRoutes(rg)

	// Tambahkan route public modul lain di sini jika ada nanti
}