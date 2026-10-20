package logger

import (
	"github.com/google/wire"
)

func ProvideZapLogger(appEnv string) (Logger, error) {
	return NewLogger(appEnv)
}

var ProviderSet = wire.NewSet(
	ProvideZapLogger,
)
var LoggerSet = wire.NewSet(
	NewLogger,
)
