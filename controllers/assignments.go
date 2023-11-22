package controllers

import (
	"csye6225-mainproject/log"
	"csye6225-mainproject/models"
	"csye6225-mainproject/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

func GetGetAllAssignmentsHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests.assignments.getall", 1)

		logger := log.GetLoggerInstance()

		if !isBodyEmpty(context) {
			logger.Debug("Request body is not empty")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Request body should be empty",
			})
			return
		}

		assignments, err := provider.MyAssignmentStore.GetAll()

		if err != nil {
			logger.Debug(fmt.Sprintf("Error getting assignments:%v", err))
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			logger.Debug(fmt.Sprintf("Assignments retrieved successfully"))
			context.JSON(http.StatusOK, assignments)
			return
		}
	}

}

func GetGetSingleAssignmentHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests.assignments.getbyid", 1)

		logger := log.GetLoggerInstance()

		if !isBodyEmpty(context) {
			logger.Debug("Request body is not empty")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Request body should be empty",
			})
			return
		}

		assignmentID := context.Param("id")

		_, err := uuid.Parse(assignmentID)

		if err != nil {
			logger.Debug(fmt.Sprintf("Given ID in params is invalid:%v", assignmentID))
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Given ID is not a valid UUID. Please check again",
			})
			return
		}

		assignment, err := provider.MyAssignmentStore.GetOne(assignmentID)

		if err != nil {
			logger.Debug(fmt.Sprintf("Given ID is not found:%v", assignmentID))
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Assignment with given ID not found. Please check again",
			})
			return
		} else {
			logger.Debug(fmt.Sprintf("Assignment with given ID is retrieved successfully"))
			context.JSON(http.StatusOK, assignment)
			return
		}

	}

}

func GetPostAssignmentHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests.assignments.post", 1)

		logger := log.GetLoggerInstance()

		account := context.MustGet("currentUserAccount").(models.Account)

		assignment, errors := convertBodyToValidAssignment(context)

		if errors != nil {
			logger.Debug("Request body is invalid")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		assignment = &models.Assignment{
			AccountID:     account.ID,
			Name:          assignment.Name,
			Points:        assignment.Points,
			NumOfAttempts: assignment.NumOfAttempts,
			Deadline:      assignment.Deadline,
		}

		updatedAssignment, err := provider.MyAssignmentStore.AddOne(assignment)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error creating the assignment:%v", err))
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			logger.Debug(fmt.Sprintf("Assignment created successfully"))
			provider.MyStatsStore.Client.Incr("assignments.count", 1)
			context.JSON(http.StatusCreated, updatedAssignment)
			return
		}
	}
}

func GetDeleteAssignmentsHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests.assignments.delete", 1)

		logger := log.GetLoggerInstance()

		if !isBodyEmpty(context) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Request body should be empty",
			})
			return
		}

		account := context.MustGet("currentUserAccount").(models.Account)

		assignmentID := context.Param("id")

		_, err := uuid.Parse(assignmentID)

		if err != nil {
			logger.Debug("Request body is not empty")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Given ID is not a valid UUID. Please check again",
			})
			return
		}

		assignment, err := provider.MyAssignmentStore.GetOne(assignmentID)

		if err != nil {
			logger.Debug(fmt.Sprintf("Assignment with given ID is not found:%v", assignmentID))
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Assignment with given ID not found. Please check again",
			})
			return
		}

		if assignment.AccountID != account.ID {
			logger.Debug(fmt.Sprintf("User %v is not authorized to delete this assingment:%v", account.ID, assignmentID))
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to delete this assignment",
			})
			return
		}

		_, err = provider.MyAssignmentStore.DeleteOne(&assignment)

		if err != nil {
			logger.Debug(fmt.Sprintf("Error deleting the assignment:%v", err))
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			provider.MyStatsStore.Client.Decr("assignments.count", 1)
			logger.Debug(fmt.Sprintf("Assignment deleted successfully"))
			context.String(http.StatusNoContent, "")
			return
		}

	}

}

func GetPutAssignmentsHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests.assignments.put", 1)

		logger := log.GetLoggerInstance()

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
			logger.Debug(fmt.Sprintf("Assignment with given ID is not found:%v", assignmentID))
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Assignment with given ID not found. Please check again",
			})
			return
		}

		if assignment.AccountID != account.ID {
			logger.Debug(fmt.Sprintf("User %v is not authorized to update this assingment:%v", account.ID, assignmentID))
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to modify this assignment",
			})
			return
		}

		assignmentFromBody, errors := convertBodyToValidAssignment(context)

		if errors != nil {
			logger.Debug("Request body is invalid")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		assignment = models.Assignment{
			Name:          assignmentFromBody.Name,
			Points:        assignmentFromBody.Points,
			NumOfAttempts: assignmentFromBody.NumOfAttempts,
			Deadline:      assignmentFromBody.Deadline,
		}

		updatedAssignment, err := provider.MyAssignmentStore.UpdateOneWithID(&assignment, assignmentID)

		if err != nil {
			logger.Debug(fmt.Sprintf("Error updating assignment:%v", err))
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			logger.Debug(fmt.Sprintf("Assignment updated successfully"))
			context.JSON(http.StatusOK, updatedAssignment)
			return
		}

	}

}

func convertBodyToValidAssignment(context *gin.Context) (*models.Assignment, []string) {
	jsonBody := extractBodyFromRequest(context.Request.Body)

	var bodyMap map[string]interface{}

	err := json.Unmarshal(jsonBody, &bodyMap)

	var validFields = map[string]string{"name": "", "num_of_attempts": "", "points": "", "deadline": "", "assignment_created": "", "assignment_updated": ""}

	if err != nil {
		return nil, []string{"Invalid JSON. Please check again"}
	}
	var errors []string
	for k, _ := range bodyMap {
		_, exists := validFields[k]
		if !exists {
			errors = append(errors, fmt.Sprintf("'%s' field in body is invalid", k))
		}
	}

	if len(errors) > 0 {
		return nil, errors
	}

	assignment := &models.Assignment{}

	if name, ok := bodyMap["name"].(string); ok {
		assignment.Name = name
	} else {
		errors = append(errors, "name should be a string")
	}

	if points, ok := bodyMap["points"].(float64); ok {

		if points == float64(int64(points)) {
			if points > 0 && points < 101 {
				assignment.Points = int(points)
			} else {
				errors = append(errors, "points should be between 1 and 100")
			}
		} else {
			errors = append(errors, "points should be a valid integer number")
		}

	} else {
		errors = append(errors, "points should be present and must a valid integer number")
	}

	if numOfAttempts, ok := bodyMap["num_of_attempts"].(float64); ok {
		if numOfAttempts == float64(int64(numOfAttempts)) {
			if numOfAttempts > 0 && numOfAttempts < 101 {
				assignment.NumOfAttempts = int(numOfAttempts)
			} else {
				errors = append(errors, "num_of_attempts should be between 1 and 100")
			}
		} else {
			errors = append(errors, "num_of_attempts should be a valid integer number")
		}
	} else {
		errors = append(errors, "num_of_attempts should be present and must be a valid integer number")
	}

	if deadline, ok := bodyMap["deadline"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, deadline)
		if err != nil {
			errors = append(errors, "deadline should be a valid date-time string. Ex: "+time.RFC3339)
		} else {
			if parsedTime.Before(time.Now()) {
				errors = append(errors, "deadline should be a valid time in the future")
			} else {
				assignment.Deadline = parsedTime
			}
		}
	} else {
		errors = append(errors, "deadline should be a valid date-time string. Ex: "+time.RFC3339)
	}

	if len(errors) > 0 {
		return nil, errors
	} else {
		return assignment, nil
	}
}

func isBodyEmpty(context *gin.Context) bool {

	contentLength := context.Request.ContentLength

	return contentLength == 0
}

func extractBodyFromRequest(readCloser io.ReadCloser) []byte {

	logger := log.GetLoggerInstance()

	body, err := io.ReadAll(readCloser)

	if err != nil {
		logger.Warn("Error reading request body")
		return make([]byte, 0)
	} else {
		return body
	}
}
