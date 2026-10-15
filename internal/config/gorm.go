package config

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGORM(sqlDB *sql.DB) (*gorm.DB, error) {
	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn: sqlDB, // Passing ko
			// neksi sql.DB yang udah ada
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			PrepareStmt: true,
			
		})

	if err != nil {
		return nil, err
	}

	fmt.Println("Initialize GORM Success")
	return gormDB, nil
}