package main

import (
	"go-starterkit/internal/config"
	"log"
)

// Alias biar gak bentrok sama 'package gen' lokal

// Execute bertugas nge-run proses generator
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[GORM Gen Error] Connection Failed config: %v", err)
	}

	sqlDB, _, err := config.NewMySQLDB(cfg)
	if err != nil {
		log.Fatalf("[GORM Gen Error] Connection Failed DB: %v", err)
	}

	gormDB, err := config.InitGORM(sqlDB)

	config.RunGenerator(
		gormDB,
		"./internal/database/query",
		"./internal/database/models",
	)
}
