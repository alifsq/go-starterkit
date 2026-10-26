package bootstrap

import (
	"context"
	"fmt"
	"go-starterkit/internal/config"
	"go-starterkit/internal/package/logger"
	"go-starterkit/internal/router"
	"log"
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Config     *config.Config
	AppHandler *router.AppHandler
	Worker     *Worker
	Enforcer   *casbin.Enforcer
	Logger     logger.Logger
}

func NewApp(cfg *config.Config, handler *router.AppHandler, worker *Worker, enforcer *casbin.Enforcer, logger logger.Logger) *App {
	return &App{
		Config:     cfg,
		AppHandler: handler,
		Worker:     worker,
		Enforcer:   enforcer,
		Logger:     logger,
	}
}

func (app *App) Start() error {
	// Set Gin Mode berdasarkan Environment
	if strings.ToLower(app.Config.AppEnv) == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Opsional: Kamu bisa pastikan policy ter-load ulang saat app start
	if err := app.Enforcer.LoadPolicy(); err != nil {
		app.Logger.Error("failed to load casbin policies on startup", zap.Error(err))
	}

	// Initialize Engine
	r := gin.Default()
	var trustedProxies []string
	for _, proxy := range strings.Split(app.Config.TrustedProxies, ",") {
		if proxy = strings.TrimSpace(proxy); proxy != "" {
			trustedProxies = append(trustedProxies, proxy)
		}
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		return fmt.Errorf("invalid TRUSTED_PROXIES configuration: %w", err)
	}

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
