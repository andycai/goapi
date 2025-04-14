package models

import (
	"time"
)

type RepoSyncRecord struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	Revision       string    `json:"revision" gorm:"size:40;not null;index"`
	Comment        string    `json:"comment" gorm:"type:text"`
	Author         string    `json:"author" gorm:"size:100"`
	SyncTime       time.Time `json:"sync_time" gorm:"not null"`
	Status         int       `json:"status" gorm:"not null;default:1"` // 1: Synced, 2: Failed
	AffectedIssues string    `json:"affected_issues" gorm:"size:255"`  // 受影响的问题列表
	AffectedFiles  string    `json:"affected_files" gorm:"type:text"`  // 受影响的文件列表，格式: A:file1.txt,M:file2.txt,D:file3.txt
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
