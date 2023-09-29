package routes

import (
	"csye6225-mainproject/db"
	"github.com/gin-gonic/gin"
)

func SetupRouter(dbHelper db.DatabaseHelper) *gin.Engine {

	r := gin.Default()

	r.Any("/healthz", createHealthzHandler(dbHelper))
	r.NoRoute(invalidHandler)

	return r
}
