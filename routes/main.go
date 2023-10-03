package routes

import (
	"csye6225-mainproject/controllers"
	"csye6225-mainproject/services"
	"csye6225-mainproject/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter(provider *services.ServiceProvider) *gin.Engine {

	r := gin.Default()

	unauthorized := r.Group("/")
	unauthorized.Any("/healthz", controllers.GetHealthzHandler(provider))

	v1 := r.Group("/v1")
	addAssignmentRoutes(v1, provider)

	r.NoRoute(utils.InvalidHandler)

	return r
}
