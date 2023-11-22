package controllers

import (
	"csye6225-mainproject/log"
	"csye6225-mainproject/models"
	"csye6225-mainproject/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

func GetPostSubmissionHandler(provider *services.ServiceProvider) func(*gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests.submissions.post", 1)

		logger := log.GetLoggerInstance()

		account := context.MustGet("currentUserAccount").(models.Account)

		assignmentID := context.Param("id")

		submission, errors := convertBodyToValidSubmission(context)

		if errors != nil {
			logger.Debug("Request body is invalid")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

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

		submission = &models.Submission{
			AccountID:         account.ID,
			SubmissionUrl:     submission.SubmissionUrl,
			AssignmentID:      assignment.ID,
			AssignmentUpdated: assignment.AssignmentUpdated,
		}

		currSubmissions, err := provider.MySubmissionStore.GetSubmissionsByAssignmentIDAndAccountID(account.ID, assignment.ID)

		if err != nil {
			logger.Debug(fmt.Sprintf("Error checking current submissions:%v", err))
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			if len(currSubmissions) >= assignment.NumOfAttempts {
				logger.Debug("Number of attempts reached")
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "You have already reached the number of attempts",
				})
				return
			}
		}

		updatedSubmission, err := provider.MySubmissionStore.AddOne(submission)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error creating the submission:%v", err))
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			logger.Debug(fmt.Sprintf("Submission created successfully"))
			_ = provider.MySubmissionStore.PublishToSNS(updatedSubmission, &account, provider.MySubmissionStore.SnsClient)
			provider.MyStatsStore.Client.Incr("submissions.count", 1)
			context.JSON(http.StatusCreated, updatedSubmission)
			return
		}
	}
}

func convertBodyToValidSubmission(context *gin.Context) (*models.Submission, []string) {
	jsonBody := extractBodyFromRequest(context.Request.Body)

	var bodyMap map[string]interface{}

	err := json.Unmarshal(jsonBody, &bodyMap)

	if err != nil {
		return nil, []string{"Invalid JSON. Please check again"}
	}

	var validFields = map[string]string{"submission_url": "", "submission_date": ""}
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

	submission := &models.Submission{}

	if submissionUrl, ok := bodyMap["submission_url"].(string); ok {
		if isValidURL(submissionUrl) {
			submission.SubmissionUrl = submissionUrl
		} else {
			errors = append(errors, "submission_url should be a valid url")
		}
	} else {
		errors = append(errors, "submission_url should be a valid url")
	}

	if len(errors) > 0 {
		return nil, errors
	} else {
		return submission, nil
	}
}

func isValidURL(inputURL string) bool {
	parsedURL, err := url.ParseRequestURI(inputURL)

	val := err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""

	return val
}
