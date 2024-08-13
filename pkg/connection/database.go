package connection

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"scylla/pkg/config"
	"scylla/pkg/engine"
	"scylla/pkg/exception"
	"time"
)

func GetDatabase(conf config.Database) *gorm.DB {

	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.User, conf.Pass, conf.Name)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	connection, err := db.DB()
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(100)
	connection.SetConnMaxLifetime(time.Second * time.Duration(300))

	engine.Instance = db

	fmt.Println("🚀 Connected Successfully to the Database")
	return db
}
