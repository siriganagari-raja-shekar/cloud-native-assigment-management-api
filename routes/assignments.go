package routes

import (
	"csye6225-mainproject/controllers"
	"csye6225-mainproject/services"
	"csye6225-mainproject/utils"
	"github.com/gin-gonic/gin"
)

func addAssignmentRoutes(rg *gin.RouterGroup, provider *services.ServiceProvider) {

	assignments := rg.Group("/assignments")

	assignments.Use(utils.UserExtractor(provider))
	assignments.GET("", controllers.GetGetAllAssignmentsHandler(provider))
	assignments.GET("/:id", controllers.GetGetSingleAssignmentHandler(provider))
	assignments.POST("", controllers.GetPostAssignmentHandler(provider))
	assignments.PUT("/:id", controllers.GetPutAssignmentsHandler(provider))
	assignments.DELETE("/:id", controllers.GetDeleteAssignmentsHandler(provider))
}
