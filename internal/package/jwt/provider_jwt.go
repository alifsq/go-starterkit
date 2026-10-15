package jwt

import "github.com/google/wire"
var JWTSet = wire.NewSet(
	NewJWTService,
)