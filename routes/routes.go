package routes

import (
	"csye6225-mainproject/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter(serviceProvider *services.ServiceProvider) *gin.Engine {

	r := gin.Default()

	r.Any("/healthz", createHealthzHandler(serviceProvider))
	r.NoRoute(invalidHandler)

	return r
}
