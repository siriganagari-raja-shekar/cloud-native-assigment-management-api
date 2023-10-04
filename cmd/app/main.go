package main

import (
	"csye6225-mainproject/conf"
	"csye6225-mainproject/routes"
	"csye6225-mainproject/services"
	"fmt"
	"os"
)

func init() {

}
func main() {

	config := &conf.Configuration{}
	config.Set()

	serviceProvider := &services.ServiceProvider{
		MyHealthzStore:    &services.HealthzStore{},
		MyAssignmentStore: services.AssignmentStore{},
		MyAccountStore:    services.AccountStore{},
	}

	serviceProvider.PopulateDBInServices()
	serviceProvider.InsertInitialUsersIntoDB()

	router := routes.SetupRouter(serviceProvider)

	serverPort := os.Getenv("SERVER_PORT")
	err := router.Run(":" + serverPort)
	if err != nil {
		fmt.Printf("Fatal server error")
	}
}
