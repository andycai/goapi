package models

import (
	"time"
)

// Page represents a static page in the database
type Page struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"size:255;not null"`
	Content   string    `gorm:"type:text"`
	Slug      string    `gorm:"size:255;uniqueIndex"`
	Status    string    `gorm:"size:20;default:'draft'"` // draft, published
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for the Page model
func (Page) TableName() string {
	return "pages"
}
