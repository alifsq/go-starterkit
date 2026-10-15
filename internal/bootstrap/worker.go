package bootstrap

import (
	"context"
	"log"
	"time"

	"go-starterkit/internal/modules/auth"
)

type Worker struct {
	authSvc auth.ServiceAuth
}

func NewWorker(authSvc auth.ServiceAuth) *Worker {
	return &Worker{authSvc: authSvc}
}

// StartTokenCleanup menjalankan ticker 24 jam untuk pembersihan DB
func (w *Worker) StartTokenCleanup(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		log.Println("[WORKER] Token cleanup worker started...")
		for {
			select {
			case <-ticker.C:
				cleanupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				if err := w.authSvc.CleanupExpiredTokens(cleanupCtx); err != nil {
					log.Printf("[WORKER ERROR] Failed Clean Up token expired: %v", err)
				} else {
					log.Println("[WORKER SUCCESS] Success Clean Up expired refresh tokens")
				}
				cancel()
			case <-ctx.Done():
				log.Println("[WORKER] Stopping token cleanup worker...")
				ticker.Stop()
				return
			}
		}
	}()
}
