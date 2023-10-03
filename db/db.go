package db

import (
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

func GetDBConf() map[string]string {
	return map[string]string{
		"host":     os.Getenv("POSTGRES_HOST"),
		"port":     os.Getenv("POSTGRES_PORT"),
		"user":     os.Getenv("POSTGRES_USER"),
		"password": os.Getenv("POSTGRES_PASSWORD"),
		"dbname":   os.Getenv("POSTGRES_DB"),
	}
}

func CreateGORMConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
}
