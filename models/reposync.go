package models

import (
	"time"
)

type RepoSyncRecord struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Revision  string    `json:"revision" gorm:"size:40;not null;index"`
	Comment   string    `json:"comment" gorm:"type:text"`
	Author    string    `json:"author" gorm:"size:100"`
	SyncTime  time.Time `json:"sync_time" gorm:"not null"`
	Status    int       `json:"status" gorm:"not null;default:1"` // 1: Synced, 2: Failed
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
