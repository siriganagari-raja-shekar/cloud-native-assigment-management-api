package services

import (
	"csye6225-mainproject/db"
	"errors"
	"gorm.io/gorm"
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

	return postgresDB.Close()
}

func (hs *HealthzStore) Ping() (bool, error) {

	if hs.db == nil {
		err := hs.OpenDBConnection(db.CreateDialectorFromEnv(), db.CreateDBConfig())
		if err != nil {
			return false, err
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
