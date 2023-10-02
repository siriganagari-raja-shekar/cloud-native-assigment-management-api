package models

import (
	"time"
)

type Account struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `json:"-"`
	AccountCreated time.Time `gorm:"autoCreateTime" json:"account_created"`
	AccountUpdated time.Time `gorm:"autoUpdateTime" json:"account_updated"`
}
