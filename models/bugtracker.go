package models

import (
	"time"
)

type Project struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	Status      int       `json:"status" gorm:"not null;default:1"` // 1: Active, 2: Completed, 3: Archived
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Iteration struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ProjectID   int64     `json:"project_id" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	Status      int       `json:"status" gorm:"not null;default:1"` // 1: Planning, 2: In Progress, 3: Completed
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Issue struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ProjectID   int64     `json:"project_id" gorm:"not null;index"`
	IterationID int64     `json:"iteration_id" gorm:"index"`
	Title       string    `json:"title" gorm:"size:200;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Type        int       `json:"type" gorm:"not null;default:1"`     // 1: Bug, 2: Feature, 3: Task
	Priority    int       `json:"priority" gorm:"not null;default:2"` // 1: Low, 2: Medium, 3: High, 4: Critical
	Status      int       `json:"status" gorm:"not null;default:1"`   // 1: New, 2: In Progress, 3: Testing, 4: Resolved, 5: Closed
	AssigneeID  int64     `json:"assignee_id" gorm:"index"`
	ReporterID  int64     `json:"reporter_id" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Comment struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	IssueID   int64     `json:"issue_id" gorm:"not null;index"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
