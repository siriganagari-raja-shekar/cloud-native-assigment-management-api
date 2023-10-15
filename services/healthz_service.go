package services

import (
	"csye6225-mainproject/db"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type HealthzStore struct {
	db *gorm.DB
}

func (hs *HealthzStore) OpenDBConnection(dialector gorm.Dialector, config *gorm.Config) error {

	if hs.db != nil {
		return nil
	}
	gormDBInstance, err := gorm.Open(dialector, config)

	if err != nil {
		hs.db = nil
		return err
	} else {
		hs.db = gormDBInstance
		return nil
	}
}

func (hs *HealthzStore) GetDBConnection() *gorm.DB {
	return hs.db
}

func (hs *HealthzStore) CloseDBConnection() error {
	if hs.db == nil {
		return errors.New("connection instance is nil")
	}

	postgresDB, err := hs.db.DB()

	if err != nil {
		return err
	}

	hs.db = nil
	return postgresDB.Close()
}

func (hs *HealthzStore) Ping() (bool, error) {

	if hs.db == nil {

		dbConf := db.GetDBConf()

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", dbConf["host"], dbConf["port"], dbConf["user"], dbConf["password"])

		err := hs.OpenDBConnection(postgres.Open(dsn), db.CreateGORMConfig())

		if err != nil {
			fmt.Printf("Unable to open DB connection with error: %v", err)
			return false, err
		}

		createDBCommand := fmt.Sprintf("CREATE DATABASE %s", os.Getenv("POSTGRES_DB"))

		res := hs.db.Exec(createDBCommand)

		if res.Error != nil {
			fmt.Printf("Database already exists: %v\n", res.Error)
		} else {
			fmt.Printf("Database created successfully\n")
		}
		err = hs.CloseDBConnection()

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf["host"], dbConf["port"], dbConf["user"], dbConf["password"], dbConf["dbname"])

		err = hs.OpenDBConnection(postgres.Open(dsn), db.CreateGORMConfig())

		if err != nil {
			return false, err
		} else {

		}
	}

	sqlDB, err := hs.db.DB()

	if err != nil {
		return false, err
	}

	pingErr := sqlDB.Ping()

	if pingErr != nil {
		return false, pingErr
	} else {
		return true, nil
	}
}
