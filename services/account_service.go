package services

import (
	"csye6225-mainproject/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountStore struct {
	Database *gorm.DB
}

func (as *AccountStore) GetAll() ([]models.Account, error) {
	var accounts []models.Account

	if err := as.Database.Find(&accounts).Error; err != nil {
		return nil, err
	} else {
		return accounts, nil
	}
}

func (as *AccountStore) GetOneByID(ID uint) (models.Account, error) {
	var account models.Account

	if err := as.Database.First(&account, models.Account{ID: ID}).Error; err != nil {
		return models.Account{}, err
	} else {
		return account, nil
	}
}

func (as *AccountStore) GetOneByEmail(email string) (models.Account, error) {
	var account models.Account

	if err := as.Database.First(&account, models.Account{Email: email}).Error; err != nil {
		return models.Account{}, err
	} else {
		return account, nil
	}
}

func (as *AccountStore) AddOne(account *models.Account) (*models.Account, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)

	if err != nil {
		return account, err
	}

	account.Password = string(hashedPassword)

	if result := as.Database.Create(&account); result.Error != nil {
		return account, result.Error
	} else {
		return account, nil
	}
}
func (as *AccountStore) UpdateOne(account *models.Account) (*models.Account, error) {

	if result := as.Database.Save(&account); result.Error != nil {
		return account, result.Error
	} else {
		return account, nil
	}
}

func (as *AccountStore) DeleteOne(account *models.Account) (*models.Account, error) {
	if result := as.Database.Delete(&account); result.Error != nil {
		return account, result.Error
	} else {
		return account, nil
	}
}
