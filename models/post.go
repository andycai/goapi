package models

import (
	"time"
)

// Post represents a blog post in the database
type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"size:255;not null"`
	Content   string    `gorm:"type:text"`
	Slug      string    `gorm:"size:255;uniqueIndex"`
	Status    string    `gorm:"size:20;default:'draft'"` // draft, published
	AuthorID  uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for the Post model
func (Post) TableName() string {
	return "posts"
}
