package auth

import "github.com/google/wire"

var AuthSet = wire.NewSet(
	NewRepositoryAuth,
	NewServiceAuth,
	NewHandlerAuth,
	NewRouterAuth,
)
