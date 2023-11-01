package main

import (
	"csye6225-mainproject/conf"
	"csye6225-mainproject/routes"
	"csye6225-mainproject/services"
	"fmt"
	"github.com/smira/go-statsd"
	"log/slog"
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
		MyStatsStore: services.StatsStore{
			Client: statsd.NewClient(os.Getenv("STATSD_SERVER_IP")+":"+os.Getenv("STATSD_SERVER_PORT"), statsd.MetricPrefix("webapp.")),
		},
	}

	serviceProvider.PopulateDBInServices()
	serviceProvider.InsertInitialUsersIntoDB()
	router := routes.SetupRouter(serviceProvider)

	serverPort := os.Getenv("SERVER_PORT")
	err := router.Run(":" + serverPort)
	if err != nil {
		slog.Error(fmt.Sprintf("Fatal server error: %v", err))
	}
}
