package services

import (
	"csye6225-mainproject/db"
	"csye6225-mainproject/log"
	"csye6225-mainproject/models"
	"encoding/csv"
	"fmt"
	"os"
)

type ServiceProvider struct {
	MyAccountStore    AccountStore
	MyAssignmentStore AssignmentStore
	MyHealthzStore    db.DatabaseHelper
	MyStatsStore      StatsStore
}

func (s *ServiceProvider) PopulateDBInServices() {
	connected, err := s.MyHealthzStore.Ping()
	logger := log.GetLoggerInstance()

	if !connected {
		logger.Error(fmt.Sprintf("Unable to connect to database and popuulate services: %s", err))
		return
	}
	s.MyAccountStore.Database = s.MyHealthzStore.GetDBConnection()
	s.MyAssignmentStore.Database = s.MyHealthzStore.GetDBConnection()
}

func (s *ServiceProvider) InsertInitialUsersIntoDB() {

	logger := log.GetLoggerInstance()
	connected, _ := s.MyHealthzStore.Ping()

	if !connected {
		return
	}

	err := s.MyAssignmentStore.Database.AutoMigrate(&models.Account{})
	if err != nil {
		logger.Error(fmt.Sprintf("Init process: Error migrating accounts: %v", err))
	}

	err = s.MyAccountStore.Database.AutoMigrate(&models.Assignment{})
	if err != nil {
		logger.Error(fmt.Sprintf("Init process: Error migrating assignments: %v", err))
	}

	logger.Info("Init process: Successfully migrated models")

	file, err := os.Open(os.Getenv("ACCOUNT_CSV_PATH"))

	if err != nil {
		logger.Error(fmt.Sprintf("Init process: Error opening file. Check file path and permissions : %v", err))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()

	if err != nil {
		logger.Error(fmt.Sprintf("Init process: Error reading lines from file: %v", err))
		return
	}

	for i := 1; i < len(lines); i++ {
		account := &models.Account{
			FirstName: lines[i][0],
			LastName:  lines[i][1],
			Email:     lines[i][2],
			Password:  lines[i][3],
		}

		account, err := s.MyAccountStore.AddOne(account)

		if err != nil {
			logger.Warn(fmt.Sprintf("Init process: Error adding user to database: %v", err))
		}
	}

	logger.Info("Init process: Successfully updated default users")

}
