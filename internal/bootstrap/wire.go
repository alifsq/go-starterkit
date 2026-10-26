//go:build wireinject
// +build wireinject

package bootstrap

import (
	"go-starterkit/internal/config"
	"go-starterkit/internal/modules/auth"
	user "go-starterkit/internal/modules/users"
	"go-starterkit/internal/package/jwt"
	"go-starterkit/internal/package/mail"
	"go-starterkit/internal/package/rbac"
	"go-starterkit/internal/package/redis"
	"go-starterkit/internal/router"

	"github.com/google/wire"
)

func InitializeApp() (*App, func(), error) {
	wire.Build(
		config.LoadConfig, // Provider Config
		wire.FieldsOf(new(*config.Config), "JwtSecret"),
		config.DatabaseSet, // Provider Database (GORM)
		jwt.JWTSet,         // Provider JWT Service <--- (PENTING)
		user.UserSet,       // Provider User Repo, Service, Handler
		auth.AuthSet,       // Provider Auth Repo, Service, Handler <--- (PENTING)
		router.RouterSet,   // Provider AppHandler
		NewWorker,          // Provider Worker Cleanup
		NewApp,             // Provider App Engine
		NewZapLogger,       // Provider Logger
		NewMailerConfig,    // Provider Mailer Config
		mail.MailSet,       // Provider Mailer
		rbac.RBACSet,       // Provider RBAC casbin
		redis.RedisSet,     // Provider Redis
	)
	return nil, nil, nil

}
