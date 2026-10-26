package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// App
	AppName        string `mapstructure:"APP_NAME"`
	AppEnv         string `mapstructure:"APP_ENV"`
	AppPort        string `mapstructure:"APP_PORT"`
	AppTimeout     string `mapstructure:"APP_TIMEOUT"`
	ApiKey         string `mapstructure:"API_KEY"`
	AppBaseURL     string `mapstructure:"APP_BASE_URL"`
	TrustedProxies string `mapstructure:"TRUSTED_PROXIES"`

	// Database
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`
	// DB Pool
	DBMaxOpenConns    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime string `mapstructure:"DB_CONN_MAX_LIFETIME"`

	// Mail Config
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     string `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`

	// Migrate DB
	MigrateDatabaseURL string `mapstructure:"MIGRATE_DATABASE_URL"`
	MigratePath        string `mapstructure:"MIGRATE_PATH"`

	// JWT Config
	JwtSecret          string        `mapstructure:"JWT_SECRET"`
	JwtAccessTokenExp  time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXP"`
	JwtRefreshTokenExp time.Duration `mapstructure:"JWT_REFRESH_TOKEN_EXP"`

	// REDIS CONFIG
	Addr         string        `mapstructure:"REDIS_ADDR"`
	Password     string        `mapstructure:"REDIS_PASSWORD"`
	DB           int           `mapstructure:"REDIS_DB"`
	PoolSize     int           `mapstructure:"REDIS_POOL_SIZE"`
	DialTimeout  time.Duration `mapstructure:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"REDIS_WRITE_TIMEOUT"`
}

var (
	appConfig *Config
	cfgOnce   sync.Once
	loadErr   error
)

func LoadConfig() (*Config, error) {

	cfgOnce.Do(func() {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			loadErr = fmt.Errorf("gagal membaca file config: %w", err)
			return // Keluar dari closure func(), bukan dari LoadConfig()
		}
		configTemp := &Config{}
		if err := viper.Unmarshal(configTemp); err != nil {
			loadErr = fmt.Errorf("gagal unmarshal config: %w", err)
			return
		}
		appConfig = configTemp
	})
	if loadErr != nil {
		return nil, loadErr
	}
	return appConfig, nil

}
