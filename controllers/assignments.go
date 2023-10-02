package controllers

import (
	"csye6225-mainproject/models"
	"github.com/gin-gonic/gin"
)

func GetAllAssignments(context *gin.Context) {
	account := context.MustGet("currentUserAccount")

	account = account.(models.Account)

}
