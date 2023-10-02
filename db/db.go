package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

type DatabaseHelper interface {
	OpenDBConnection(dialector gorm.Dialector, config *gorm.Config) error
	GetDBConnection() *gorm.DB
	CloseDBConnection() error
	Ping() (bool, error)
}

func CreateDialectorFromEnv() gorm.Dialector {
	conf := map[string]string{
		"host":     os.Getenv("POSTGRES_HOST"),
		"port":     os.Getenv("POSTGRES_PORT"),
		"user":     os.Getenv("POSTGRES_USER"),
		"password": os.Getenv("POSTGRES_PASSWORD"),
		"dbname":   os.Getenv("POSTGRES_DB"),
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf["host"], conf["port"], conf["user"], conf["password"], conf["dbname"])
	return postgres.Open(dsn)
}

func CreateDBConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
}
