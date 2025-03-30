package models

import "time"

type PatchRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OldVersion  string    `gorm:"size:50" json:"old_version"`
	NewVersion  string    `gorm:"size:50" json:"new_version"`
	Version     string    `gorm:"size:50;not null" json:"version"`
	PatchFile   string    `gorm:"size:255;not null" json:"patch_file"`
	Size        int64     `gorm:"not null" json:"size"`
	FileCount   int       `gorm:"default:0" json:"file_count"`
	Status      int       `gorm:"default:0" json:"status"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
