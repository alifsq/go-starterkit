package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestUnmarshalRedisConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigType("env")
	if err := v.ReadConfig(strings.NewReader(`
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=test-password
REDIS_DB=2
REDIS_POOL_SIZE=12
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=4s
TRUSTED_PROXIES=10.0.0.0/8,192.168.0.0/16
`)); err != nil {
		t.Fatal(err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.Addr != "localhost:6379" {
		t.Errorf("expected Redis address localhost:6379, got %q", cfg.Addr)
	}
	if cfg.Password != "test-password" {
		t.Errorf("expected Redis password to be unmarshaled")
	}
	if cfg.DB != 2 {
		t.Errorf("expected Redis DB 2, got %d", cfg.DB)
	}
	if cfg.PoolSize != 12 {
		t.Errorf("expected Redis pool size 12, got %d", cfg.PoolSize)
	}
	if cfg.DialTimeout.String() != "5s" {
		t.Errorf("expected Redis dial timeout 5s, got %s", cfg.DialTimeout)
	}
	if cfg.ReadTimeout.String() != "3s" {
		t.Errorf("expected Redis read timeout 3s, got %s", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout.String() != "4s" {
		t.Errorf("expected Redis write timeout 4s, got %s", cfg.WriteTimeout)
	}
	if cfg.TrustedProxies != "10.0.0.0/8,192.168.0.0/16" {
		t.Errorf("expected trusted proxies to be unmarshaled, got %q", cfg.TrustedProxies)
	}
}
