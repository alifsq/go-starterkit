package bootstrap

import (
	"context"
	"go-starterkit/internal/config"
	"go-starterkit/internal/router"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config     *config.Config
	AppHandler *router.AppHandler
	Worker     *Worker
}

func NewApp(cfg *config.Config, handler *router.AppHandler, worker *Worker) *App {
	return &App{
		Config:     cfg,
		AppHandler: handler,
		Worker:     worker,
	}
}

func (app *App) Start() error {
	// Set Gin Mode berdasarkan Environment
	if strings.ToLower(app.Config.AppEnv) == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Engine
	r := gin.Default()

	// Register Middleware & Routes
	app.AppHandler.RegisterRoutes(r)

	// Context untuk lifecycle Background Worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// BackGround Worker (Token Clean Up)
	if app.Worker != nil {
		app.Worker.StartTokenCleanup(ctx)
	}

	// Format port (:8080)
	port := app.Config.AppPort
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("Application [%s] is running on port %s in %s mode", app.Config.AppName, port, app.Config.AppEnv)

	// Graceful Shutdown
	srv := NewServer(port, r)
	return srv.Run()
}
