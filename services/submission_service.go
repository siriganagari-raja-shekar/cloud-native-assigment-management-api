package services

import (
	"csye6225-mainproject/log"
	"csye6225-mainproject/models"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"gorm.io/gorm"
	"os"
	"strings"
)

type SubmissionStore struct {
	Database  *gorm.DB
	SnsClient *sns.SNS
}

func (ss *SubmissionStore) AddOne(submission *models.Submission) (*models.Submission, error) {
	if err := ss.Database.Create(&submission).Error; err != nil {
		return nil, err
	} else {
		return submission, nil
	}
}

func (ss *SubmissionStore) GetSubmissionsByAssignmentIDAndAccountID(accountID string, assignmentID string) ([]models.Submission, error) {
	var submissions []models.Submission
	if err := ss.Database.Find(&submissions, models.Submission{AccountID: accountID, AssignmentID: assignmentID}).Error; err != nil {
		return nil, err
	} else {
		return submissions, nil
	}
}

func (ss *SubmissionStore) PublishToSNS(submission *models.Submission, account *models.Account, client *sns.SNS) error {
	logger := log.GetLoggerInstance()
	type Message struct {
		AssignmentID  string `json:"assignment_id"`
		AccountID     string `json:"account_id"`
		SubmissionUrl string `json:"submission_url"`
		Email         string `json:"email"`
	}

	message := Message{
		AssignmentID:  submission.AssignmentID,
		AccountID:     submission.AccountID,
		SubmissionUrl: submission.SubmissionUrl,
		Email:         account.Email,
	}

	messageMarshalled, err := json.MarshalIndent(message, "", strings.Repeat(" ", 4))
	if err != nil {
		logger.Error(fmt.Sprintf("Error marshalling message: %v", err))
		os.Exit(1)
	} else {
		logger.Debug(fmt.Sprintf("Message marshalled: %v", string(messageMarshalled)))
	}

	result, err := client.Publish(&sns.PublishInput{
		Message:  aws.String(string(messageMarshalled)),
		TopicArn: aws.String(os.Getenv("SNS_TOPIC_ARN")),
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Error sending message to SNS: %v", err))
	} else {
		logger.Debug(fmt.Sprintf("Result: %v", result.String()))
	}

	return err

}
