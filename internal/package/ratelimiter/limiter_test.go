package ratelimiter

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestCheckLimitRejectsUnknownPolicy(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	defer client.Close()
	engine := NewLimiterEngine(client)

	if _, err := engine.CheckLimit(context.Background(), "not-configured", "client"); err == nil {
		t.Fatal("expected unknown rate limit policy to return an error")
	}
}

func TestCheckLimitRejectsEmptyIdentifier(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	defer client.Close()
	engine := NewLimiterEngine(client)

	if _, err := engine.CheckLimit(context.Background(), "login_ip", ""); err == nil {
		t.Fatal("expected empty identifier to return an error")
	}
}
