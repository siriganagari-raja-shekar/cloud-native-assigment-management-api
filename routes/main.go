package routes

import (
	"csye6225-mainproject/controllers"
	"csye6225-mainproject/services"
	"csye6225-mainproject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

func SetupRouter(provider *services.ServiceProvider) *gin.Engine {

	file, err := os.OpenFile(os.Getenv("REQUEST_LOG_FILE_PATH"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error opening the file: %v", err))
	} else {
		gin.DefaultWriter = file
	}

	r := gin.Default()

	r.Use(utils.StatsLogger(provider))

	unauthorized := r.Group("/")
	unauthorized.Any("/healthz", controllers.GetHealthzHandler(provider))

	v1 := r.Group("/v1")
	addAssignmentRoutes(v1, provider)

	r.NoRoute(utils.InvalidHandler(provider))

	return r
}
