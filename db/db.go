package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

type DatabaseHelper interface {
	OpenDBConnection(dialector gorm.Dialector, config *gorm.Config)
	GetDBConnection() *gorm.DB
	CloseDBConnection() error
}

type Store struct {
	db *gorm.DB
}

func (s *Store) OpenDBConnection(dialector gorm.Dialector, config *gorm.Config) {

	gormDBInstance, err := gorm.Open(dialector, config)

	if err != nil {
		s.db = nil
	} else {
		s.db = gormDBInstance
	}

}

func (s *Store) GetDBConnection() *gorm.DB {
	return s.db
}

func (s *Store) CloseDBConnection() error {
	if s.db == nil {
		return errors.New("connection instance is nil")
	}

	postgresDB, err := s.db.DB()

	if err != nil {
		return err
	}

	return postgresDB.Close()
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
