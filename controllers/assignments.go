package controllers

import (
	"csye6225-mainproject/models"
	"csye6225-mainproject/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func GetGetAllAssignmentsHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {
		account := context.MustGet("currentUserAccount").(models.Account)

		assignments, err := provider.MyAssignmentStore.GetAllByAccount(account.ID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			context.JSON(http.StatusOK, assignments)
			return
		}
	}

}

func GetGetSingleAssignmentHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {

		account := context.MustGet("currentUserAccount").(models.Account)

		assignmentID := context.Param("id")

		_, err := uuid.Parse(assignmentID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Given ID is not valid. Please check again",
			})
			return
		}

		assignment, err := provider.MyAssignmentStore.GetOne(assignmentID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Assignment with given ID not found. Please check again",
			})
			return
		}

		if assignment.AccountID != account.ID {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to access this assignment",
			})
			return
		} else {
			context.JSON(http.StatusOK, assignment)
			return
		}

	}

}

func GetPostAssignmentHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {
		account := context.MustGet("currentUserAccount").(models.Account)

		assignment := models.Assignment{}

		if err := context.ShouldBindJSON(&assignment); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "One of the constraints is not satisfied",
				"constraints": []string{
					"1. Number of points should be between 1 and 100",
					"2. Number of attempts should be between 1 and 100",
					"3. Deadline should be a future date",
				},
			})
			return
		}

		if assignment.Deadline.Before(time.Now()) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Deadline should be a future date",
			})
			return
		}

		_, err := provider.MyAssignmentStore.GetOneWithAccountIDAndName(account.ID, assignment.Name)

		if err == nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Assignment with given name was already created by you",
			})
			return
		}

		assignment = models.Assignment{
			AccountID:     account.ID,
			Name:          assignment.Name,
			Points:        assignment.Points,
			NumOfAttempts: assignment.NumOfAttempts,
			Deadline:      assignment.Deadline,
		}

		updatedAssignment, err := provider.MyAssignmentStore.AddOne(&assignment)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			context.JSON(http.StatusCreated, updatedAssignment)
			return
		}
	}
}

func GetDeleteAssignmentsHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {
		account := context.MustGet("currentUserAccount").(models.Account)

		assignmentID := context.Param("id")

		_, err := uuid.Parse(assignmentID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Given ID is not valid. Please check again",
			})
			return
		}

		assignment, err := provider.MyAssignmentStore.GetOne(assignmentID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Assignment with given ID not found. Please check again",
			})
			return
		}

		if assignment.AccountID != account.ID {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to delete this assignment",
			})
			return
		}

		_, err = provider.MyAssignmentStore.DeleteOne(&assignment)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			context.String(http.StatusNoContent, "")
			return
		}

	}

}

func GetPutAssignmentsHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {
		account := context.MustGet("currentUserAccount").(models.Account)

		assignmentID := context.Param("id")

		_, err := uuid.Parse(assignmentID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Given ID is not valid. Please check again",
			})
			return
		}

		assignment, err := provider.MyAssignmentStore.GetOne(assignmentID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Assignment with given ID not found. Please check again",
			})
			return
		}

		if assignment.AccountID != account.ID {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to modify this assignment",
			})
			return
		}

		assignmentFromJSON := models.Assignment{}

		if err := context.ShouldBindJSON(&assignmentFromJSON); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "One of the constraints is not satisfied",
				"constraints": []string{
					"1. Number of points should be between 1 and 100",
					"2. Number of attempts should be between 1 and 100",
					"3. Deadline should be a future date",
					"4. An ID should not be specified in the request body",
				},
			})
			return
		}

		if len(strings.TrimSpace(assignmentFromJSON.ID)) > 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Don't specify an ID in the request body",
			})
			return
		}

		if assignmentFromJSON.Deadline.Before(time.Now()) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Deadline should be a future date",
			})
			return
		}

		dupAssignment, err := provider.MyAssignmentStore.GetOneWithAccountIDAndName(account.ID, assignment.Name)

		if err == nil && dupAssignment.ID != assignment.ID {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Assignment with given name was already created by you",
			})
			return
		}

		assignment = models.Assignment{
			Name:          assignment.Name,
			Points:        assignment.Points,
			NumOfAttempts: assignment.NumOfAttempts,
			Deadline:      assignment.Deadline,
		}

		_, err = provider.MyAssignmentStore.UpdateOneWithID(&assignment, assignmentID)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			context.String(http.StatusNoContent, "")
			return
		}

	}

}
