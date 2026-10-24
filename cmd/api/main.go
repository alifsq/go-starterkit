package main

import (
	"go-starterkit/internal/bootstrap"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}
	app, cleaup, err := bootstrap.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	defer cleaup()

	// 2. Nyalakan server Gin
	if err := app.Start(); err != nil {
		logger.Fatal("Server stopped with error", zap.Error(err))
	}

}
