package models

import "time"

type Submission struct {
	ID                string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SubmissionUrl     string     `json:"submission_url"`
	AccountID         string     `json:"-"`
	Account           Account    `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	AssignmentID      string     `json:"assignment_id"`
	Assignment        Assignment `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	SubmissionDate    time.Time  `gorm:"autoCreateTime" json:"submission_created"`
	AssignmentUpdated time.Time  `json:"assignment_updated"`
}
