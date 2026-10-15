package bootstrap

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler *gin.Engine) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Run lifecycle startup & graceful shutdown
func (s *Server) Run() error {
serverErr := make(chan error, 1)

	// 1. Run HTTP Server
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// 2. Listen OS Signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err // <--If the port is already in use / other startup error
	case sig := <-quit:
		log.Printf("Received signal [%s]. Shutting down gracefully...", sig)
	}

	// 3. Context Timeout For Graceful Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop HTTP Server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err 
	}

	log.Println("Server exited gracefully")
	return nil // Return NIL agar exit code program adalah 0!
}