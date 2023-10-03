package main

import (
	"csye6225-mainproject/conf"
	"csye6225-mainproject/routes"
	"csye6225-mainproject/services"
	"fmt"
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

	err := router.Run(":8000")
	if err != nil {
		fmt.Printf("Fatal server error")
	}
}
