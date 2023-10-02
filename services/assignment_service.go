package services

import (
	"csye6225-mainproject/models"
	"gorm.io/gorm"
)

type AssignmentStore struct {
	Database *gorm.DB
}

func (as *AssignmentStore) GetAllByUser(accountID uint) ([]models.Assignment, error) {
	var assignments []models.Assignment

	if err := as.Database.Find(&assignments, models.Assignment{AccountID: accountID}).Error; err != nil {
		return nil, err
	} else {
		return assignments, nil
	}
}

func (as *AssignmentStore) GetOne(ID uint) (models.Assignment, error) {
	var assignment models.Assignment

	if err := as.Database.First(&assignment, ID).Error; err != nil {
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
func (as *AssignmentStore) UpdateOne(assignment *models.Assignment) (*models.Assignment, error) {

	if result := as.Database.Save(&assignment); result.Error != nil {
		return assignment, result.Error
	} else {
		return assignment, nil
	}
}

func (as *AssignmentStore) DeleteOne(assignment *models.Assignment) (*models.Assignment, error) {
	if result := as.Database.Delete(&assignment); result.Error != nil {
		return assignment, result.Error
	} else {
		return assignment, nil
	}
}
