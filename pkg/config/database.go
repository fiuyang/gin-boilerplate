package config

import (
	"fmt"
	"gorm.io/gorm/logger"
	"scylla/pkg/exception"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionDB(config *Config) *gorm.DB {

	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUsername, config.DBPassword, config.DBName)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	connection, err := db.DB()
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(100)
	connection.SetConnMaxLifetime(time.Second * time.Duration(300))

	fmt.Println("🚀 Connected Successfully to the Database")
	return db
}
