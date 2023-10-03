package services

import (
	"csye6225-mainproject/models"
	"gorm.io/gorm"
)

type AssignmentStore struct {
	Database *gorm.DB
}

func (as *AssignmentStore) GetAllByAccount(accountID string) ([]models.Assignment, error) {
	var assignments []models.Assignment

	if err := as.Database.Find(&assignments, models.Assignment{AccountID: accountID}).Error; err != nil {
		return nil, err
	} else {
		return assignments, nil
	}
}

func (as *AssignmentStore) GetOne(ID string) (models.Assignment, error) {
	var assignment models.Assignment

	if err := as.Database.First(&assignment, models.Assignment{ID: ID}).Error; err != nil {
		return models.Assignment{}, err
	} else {
		return assignment, nil
	}
}

func (as *AssignmentStore) GetOneWithAccountIDAndName(ID string, name string) (models.Assignment, error) {
	var assignment models.Assignment

	if err := as.Database.First(&assignment, models.Assignment{AccountID: ID, Name: name}).Error; err != nil {
		return models.Assignment{}, err
	} else {
		return assignment, nil
	}
}

func (as *AssignmentStore) GetOneWithAccountIDAndAssignmentID(accountID string, assignmentID string) (models.Assignment, error) {
	var assignment models.Assignment

	if err := as.Database.First(&assignment, models.Assignment{AccountID: accountID, ID: assignmentID}).Error; err != nil {
		return models.Assignment{}, err
	} else {
		return assignment, nil
	}
}

func (as *AssignmentStore) AddOne(assignment *models.Assignment) (*models.Assignment, error) {
	if result := as.Database.Create(&assignment); result.Error != nil {
		return assignment, result.Error
	} else {
		return assignment, nil
	}
}

func (as *AssignmentStore) UpdateOneWithID(assignment *models.Assignment, assignmentID string) (*models.Assignment, error) {

	tempAssignment, _ := as.GetOne(assignmentID)

	tempAssignment.Name = assignment.Name
	tempAssignment.Points = assignment.Points
	tempAssignment.NumOfAttempts = assignment.NumOfAttempts
	tempAssignment.Deadline = assignment.Deadline

	if result := as.Database.Save(&tempAssignment); result.Error != nil {
		return &tempAssignment, result.Error
	} else {
		return &tempAssignment, nil
	}
}

func (as *AssignmentStore) DeleteOne(assignment *models.Assignment) (*models.Assignment, error) {
	if result := as.Database.Delete(&assignment); result.Error != nil {
		return assignment, result.Error
	} else {
		return assignment, nil
	}
}
