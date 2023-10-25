package services

import (
	"csye6225-mainproject/db"
	"csye6225-mainproject/models"
	"encoding/csv"
	"fmt"
	"os"
)

type ServiceProvider struct {
	MyAccountStore    AccountStore
	MyAssignmentStore AssignmentStore
	MyHealthzStore    db.DatabaseHelper
}

func (s *ServiceProvider) PopulateDBInServices() {
	connected, _ := s.MyHealthzStore.Ping()

	if !connected {
		return
	}
	s.MyAccountStore.Database = s.MyHealthzStore.GetDBConnection()
	s.MyAssignmentStore.Database = s.MyHealthzStore.GetDBConnection()
}

func (s *ServiceProvider) InsertInitialUsersIntoDB() {

	connected, _ := s.MyHealthzStore.Ping()

	if !connected {
		return
	}

	err := s.MyAssignmentStore.Database.AutoMigrate(&models.Account{})
	if err != nil {
		fmt.Printf("Error migrating accounts: %v\n", err)
	}

	err = s.MyAccountStore.Database.AutoMigrate(&models.Assignment{})
	if err != nil {
		fmt.Printf("Error migrating assignments\n: %v", err)
	}

	file, err := os.Open(os.Getenv("ACCOUNT_CSV_PATH"))

	if err != nil {
		fmt.Printf("Error opening file. Check file path and permissions : %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()

	if err != nil {
		fmt.Printf("Error reading lines from file: %v\n", err)
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
			fmt.Printf("Error adding user to database: %v\n", err)
		}
	}

}
